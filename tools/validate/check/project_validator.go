package check

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Josepavese/asdp/engine/domain"
)

type ValidateProjectUseCase struct {
	fs           domain.FileSystem
	hasher       domain.ContentHasher
	parser       domain.ASTParser
	configLoader domain.ConfigurationLoader
	baseConfig   *domain.Config
}

func NewValidateProjectUseCase(fs domain.FileSystem, parser domain.ASTParser, hasher domain.ContentHasher, configLoader domain.ConfigurationLoader, baseConfig *domain.Config) *ValidateProjectUseCase {
	return &ValidateProjectUseCase{
		fs:           fs,
		hasher:       hasher,
		parser:       parser,
		configLoader: configLoader,
		baseConfig:   baseConfig,
	}
}

type ValidationReport struct {
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
	IsValid  bool                `json:"is_valid"`
}

type ValidationError struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type ValidationWarning struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

func (uc *ValidateProjectUseCase) Execute(rootPath string) (*ValidationReport, error) {
	report := &ValidationReport{
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
		IsValid:  true,
	}

	// 0. Load Project Config
	config, err := uc.configLoader.LoadForProject(uc.baseConfig, rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	// 1. Root Check (Mandatory Files)
	for _, filename := range config.Validation.MandatoryFiles {
		filePath := filepath.Join(rootPath, filename)
		if _, err := uc.fs.Stat(filePath); err != nil {
			report.Errors = append(report.Errors, ValidationError{
				Path:   rootPath,
				Reason: fmt.Sprintf("Missing required file: %s at project root", filename),
			})
		}
	}

	// 2. Walk Tree
	err = uc.fs.Walk(rootPath, func(path string, isDir bool) error {
		if !isDir {
			return nil
		}
		if uc.shouldIgnoreDir(path) {
			return nil // Skip .git, .agent, etc.
		}

		// Analyze folder "significance"
		isSignificant, _, isLeaf := uc.analyzeFolderSignificance(path, config.Validation.Freshness)
		if !isSignificant {
			return nil
		}
		// If significant (HUB or LEAF), it MUST have required module files

		for _, filename := range config.Validation.ModuleFiles {
			// Special handling: codemodel is only required if it's a Leaf
			if filename == "codemodel.md" && !isLeaf {
				continue
			}

			fullPath := filepath.Join(path, filename)
			if !uc.fileExists(fullPath) {
				report.Errors = append(report.Errors, ValidationError{
					Path:   path,
					Reason: fmt.Sprintf("Missing required file: %s (Significant Module)", filename),
				})
			}
		}

		// B. Content Check (Specific to codespec currently)
		specPath := filepath.Join(path, "codespec.md")
		if uc.fileExists(specPath) {
			specContentBytes, err := uc.fs.ReadFile(specPath)
			if err == nil {
				specContent := string(specContentBytes)
				if err := uc.validateSpecContent(path, specContent, config.Validation); err != nil {
					report.Errors = append(report.Errors, ValidationError{
						Path:   path,
						Reason: err.Error(),
					})
				}

				// C. Freshness Check (Warning)
				uc.checkFreshness(path, specPath, filepath.Join(path, "codemodel.md"), config.Validation.Freshness, report)
			}
		}

		return nil
	})

	if len(report.Errors) > 0 {
		report.IsValid = false
	}

	return report, err
}

func (uc *ValidateProjectUseCase) fileExists(path string) bool {
	_, err := uc.fs.Stat(path)
	return err == nil
}

func (uc *ValidateProjectUseCase) shouldIgnoreDir(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules"
}

func (uc *ValidateProjectUseCase) validateSpecContent(path, content string, config domain.ValidationConfig) error {
	for _, forbidden := range config.ForbiddenStrings {
		if strings.Contains(content, forbidden) {
			return fmt.Errorf("codespec.md contains forbidden '%s' placeholders", forbidden)
		}
	}

	for _, required := range config.RequiredSpecKeys {
		if !strings.Contains(content, required) {
			return fmt.Errorf("codespec.md invalid structure (missing '%s')", required)
		}
	}
	return nil
}

func (uc *ValidateProjectUseCase) checkFreshness(dirPath, specPath, modelPath string, freshness domain.FreshnessConfig, report *ValidationReport) {
	// 1. Get latest Code modification time in dir
	files, err := uc.fs.ReadDir(dirPath)
	if err != nil {
		return
	}

	var maxCodeTime time.Time
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()

		isWatched := false
		for _, ext := range freshness.WatchedExtensions {
			if strings.HasSuffix(name, ext) {
				isWatched = true
				break
			}
		}

		isIgnored := false
		for _, ext := range freshness.IgnoredExtensions {
			if strings.HasSuffix(name, ext) {
				isIgnored = true
				break
			}
		}

		if isWatched && !isIgnored {
			if f.ModTime().After(maxCodeTime) {
				maxCodeTime = f.ModTime()
			}
		}
	}

	if maxCodeTime.IsZero() {
		return // No code to compare against
	}

	// 2. Check CodeSpec Freshness
	specInfo, err := uc.fs.Stat(specPath)
	if err == nil {
		if specInfo.ModTime().Before(maxCodeTime) {
			report.Warnings = append(report.Warnings, ValidationWarning{
				Path:   dirPath,
				Reason: fmt.Sprintf("Stale CodeSpec: codespec.md (%v) is older than source code (%v)", specInfo.ModTime().Format(time.RFC3339), maxCodeTime.Format(time.RFC3339)),
			})
		}
	}

	// 3. Check CodeModel Freshness
	modelInfo, err := uc.fs.Stat(modelPath)
	if err == nil {
		if modelInfo.ModTime().Before(maxCodeTime) {
			report.Warnings = append(report.Warnings, ValidationWarning{
				Path:   dirPath,
				Reason: fmt.Sprintf("Stale CodeModel: codemodel.md (%v) is older than source code (%v)", modelInfo.ModTime().Format(time.RFC3339), maxCodeTime.Format(time.RFC3339)),
			})
		}
	}
}

func (uc *ValidateProjectUseCase) analyzeFolderSignificance(path string, freshness domain.FreshnessConfig) (isSignificant bool, isHub bool, isLeaf bool) {
	files, err := uc.fs.ReadDir(path)
	if err != nil {
		return false, false, false
	}

	hasCode := false
	subDirs := 0

	for _, f := range files {
		if f.IsDir() {
			if !strings.HasPrefix(f.Name(), ".") {
				subDirs++
			}
		} else {
			name := f.Name()
			for _, ext := range freshness.WatchedExtensions {
				if strings.HasSuffix(name, ext) {
					hasCode = true
					break
				}
			}
		}
	}

	isLeaf = hasCode
	isHub = subDirs > 1

	isSignificant = isLeaf || isHub
	return isSignificant, isHub, isLeaf
}
