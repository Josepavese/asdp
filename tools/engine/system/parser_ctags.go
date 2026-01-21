package system

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
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

	// Run ctags (NON-RECURSIVE)
	// --output-format=json : standardized JSON output
	// --fields=+nK : line number, kind
	// --exclude=*.go : let GoASTParser handle Go files
	cmd := exec.Command(p.BinaryPath, "--output-format=json", "--fields=+nK", "--exclude=*.go", root)

	// We capture stdout
	var out bytes.Buffer
	cmd.Stdout = &out
	// We ignore stderr for now or log it

	if err := cmd.Run(); err != nil {
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
