package usecase

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Josepavese/asdp/engine/domain"
)

type ScaffoldUseCase struct {
	fs     domain.FileSystem
	config domain.ScaffoldConfig
}

type ScaffoldParams struct {
	Name    string
	Type    string // library, service, app
	Path    string // Parent directory, default "."
	Title   string // Required
	Summary string // Required
	Context string // Required
}

func NewScaffoldUseCase(fs domain.FileSystem, config domain.ScaffoldConfig) *ScaffoldUseCase {
	return &ScaffoldUseCase{
		fs:     fs,
		config: config,
	}
}

func (uc *ScaffoldUseCase) Execute(params ScaffoldParams) (string, error) {
	absPath, err := validateAndExpandPath(params.Path)
	if err != nil {
		return "", err
	}
	params.Path = absPath

	targetDir := params.Path
	moduleName := params.Name

	if params.Name == "." || params.Name == "" {
		targetDir = params.Path
		moduleName = filepath.Base(params.Path)
	} else {
		targetDir = filepath.Join(params.Path, params.Name)
	}

	// Validate params based on Config
	if uc.config.RequiredContext {
		if params.Title == "" || params.Summary == "" || params.Context == "" {
			return "", fmt.Errorf("Title, Summary, and Context are required for strict scaffolding")
		}
	}

	// Create Directory
	if err := uc.fs.MkdirAll(targetDir); err != nil {
		// Proceed if error is just "exists", but MkdirAll usually handles that.
	}

	// Generate Content
	// Use Templates from Config
	specContent, err := uc.renderTemplate(uc.config.SpecTemplate, map[string]interface{}{
		"Name":         moduleName,
		"Type":         params.Type,
		"ASDPVersion":  domain.Version,
		"Title":        params.Title,
		"Summary":      params.Summary,
		"Context":      params.Context,
		"LastModified": time.Now().Format(time.RFC3339),
		"ID":           moduleName, // Simple ID for now
	})
	if err != nil {
		return "", fmt.Errorf("failed to render spec template: %w", err)
	}

	modelContent, err := uc.renderTemplate(uc.config.ModelTemplate, map[string]interface{}{
		"CheckedAt":   time.Now().Format(time.RFC3339),
		"ASDPVersion": domain.Version,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render model template: %w", err)
	}

	// Write Files (Safely)
	files := map[string]string{
		"codespec.md":  specContent,
		"codemodel.md": modelContent,
	}

	created := []string{}
	skipped := []string{}

	for filename, content := range files {
		filePath := filepath.Join(targetDir, filename)
		if _, err := uc.fs.Stat(filePath); err == nil {
			skipped = append(skipped, filename)
			continue
		}

		// Attempt write
		if err := uc.fs.WriteFile(filePath, []byte(content)); err != nil {
			return "", fmt.Errorf("failed to write %s: %w", filename, err)
		}
		created = append(created, filename)
	}

	msg := fmt.Sprintf("Scaffolded %s in %s. Created: %v, Skipped: %v", moduleName, targetDir, created, skipped)

	if len(created) > 0 {
		msg += fmt.Sprintf("\n\n[SUCCESS] Module '%s' created with strict compliance.", moduleName)
	}

	return msg, nil
}

func (uc *ScaffoldUseCase) renderTemplate(tmplStr string, data interface{}) (string, error) {
	t, err := template.New("scaffold").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
