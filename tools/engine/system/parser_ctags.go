package system

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Josepavese/asdp/engine/domain"
)

type CtagsParser struct {
	// Optional: path to ctags binary, defaults to "ctags"
	BinaryPath string
}

func NewCtagsParser() *CtagsParser {
	return &CtagsParser{BinaryPath: "ctags"}
}

// CtagsEntry matches the JSON output of `ctags --output-format=json`
type CtagsEntry struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Kind      string `json:"kind"`
	Scope     string `json:"scope"`
	ScopeKind string `json:"scopeKind"`
	Line      int    `json:"line"`
	End       int    `json:"end"`
	Signature string `json:"signature"` // Custom field often not present by default
	Pattern   string `json:"pattern"`   // Regex pattern to find the line
}

func (p *CtagsParser) ParseDir(root string) ([]domain.Symbol, error) {
	// Check if ctags is available
	if _, err := exec.LookPath(p.BinaryPath); err != nil {
		// Log warning? For now just return empty, assuming optional.
		// Alternatively, return error so the Polyglot knows it failed.
		// Let's return a specific error that can be ignored.
		return nil, fmt.Errorf("ctags binary not found: %w", err)
	}

	// 1. Collect files RECURSIVELY (Boundary-Aware)
	var filesToScan []string

	ignoredPatterns := []string{
		"node_modules", "vendor", "packages", "bower_components",
		"venv", ".venv", "anaconda", "conda", "env", ".env",
		"bin", "obj", "dist", "target", "build", "out",
		".vscode", ".idea", ".git", ".hg", ".svn", ".cache",
	}

	isIgnored := func(name string) bool {
		name = strings.ToLower(name)
		for _, p := range ignoredPatterns {
			if strings.Contains(name, p) {
				return true
			}
		}
		return false
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == root {
				return nil
			}
			if isIgnored(d.Name()) {
				return filepath.SkipDir
			}
			// Boundary Check
			if _, err := os.Stat(filepath.Join(path, "codespec.md")); err == nil {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(path, "codemodel.md")); err == nil {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip Go files (handled by GoASTParser) and others
		if strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		info, err := d.Info()
		if err == nil && info.Mode().IsRegular() {
			filesToScan = append(filesToScan, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk dir: %w", err)
	}

	if len(filesToScan) == 0 {
		return []domain.Symbol{}, nil
	}

	// 2. Run ctags with file list input
	// -L - : read file list from stdin
	cmd := exec.Command(p.BinaryPath, "--output-format=json", "--fields=+nKe", "--exclude=*.go", "-L", "-")

	var out bytes.Buffer
	cmd.Stdout = &out

	// Create stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("ctags start failed: %w", err)
	}

	// Write files to stdin
	go func() {
		defer stdin.Close()
		for _, f := range filesToScan {
			fmt.Fprintln(stdin, f)
		}
	}()

	if err := cmd.Wait(); err != nil {
		// Verify if it's an exit code issue or IO
		// ctags might exit with non-zero if warnings?
		// For now simple error return
		return nil, fmt.Errorf("ctags execution failed: %w", err)
	}

	var symbols []domain.Symbol
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Bytes()
		var entry CtagsEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			continue // skip malformed lines
		}

		// Filter out "ptag" or other internal metadata if needed.
		// Universal ctags JSON has `_type` field but we mapped minimal fields.
		// Usually real tags have "name".

		if entry.Name == "" {
			continue
		}

		// Map to Domain Symbol
		sym := domain.Symbol{
			Name:      entry.Name,
			Kind:      entry.Kind,
			Line:      entry.Line,
			LineEnd:   entry.End,
			Exported:  true,          // Assume exported by default for non-Go langs
			Signature: entry.Pattern, // Use pattern as rough signature surrogate
		}

		// Improve signature if pattern is just a search regex
		// Often pattern looks like: "/^func MyFunc() {$/"
		if strings.HasPrefix(sym.Signature, "/^") {
			sym.Signature = strings.TrimSuffix(strings.TrimPrefix(sym.Signature, "/^"), "$/")
		}

		if entry.Scope != "" {
			sym.Parent = entry.Scope
		}

		symbols = append(symbols, sym)
	}

	return symbols, nil
}
