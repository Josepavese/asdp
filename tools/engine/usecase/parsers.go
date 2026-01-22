package usecase

import (
	"fmt"
	"strings"

	"github.com/Josepavese/asdp/engine/domain"
	"gopkg.in/yaml.v3"
)

// Shared helpers to parse Frontmatter

func parseCodeSpec(data []byte) (*domain.CodeSpec, error) {
	// Simple split by "---"
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return nil, nil // Invalid format
	}

	fontmatter := parts[1]
	body := parts[2]

	var meta domain.CodeSpecMeta
	if err := yaml.Unmarshal([]byte(fontmatter), &meta); err != nil {
		return nil, err
	}

	return &domain.CodeSpec{
		MetaData: meta,
		Body:     body,
	}, nil
}

// ValidateCodeSpec checks for lazy values and returns a ValidationResult.
func ValidateCodeSpec(spec *domain.CodeSpec) *domain.ValidationResult {
	if spec == nil {
		return nil
	}

	meta := spec.MetaData
	var errors []string

	// Helper to check for lazy patterns
	isLazy := func(val string) bool {
		val = strings.TrimSpace(strings.ToLower(val))
		return val == "" ||
			val == "<no value>" ||
			val == "unknown" ||
			val == "todo" ||
			strings.HasPrefix(val, "todo") ||
			strings.Contains(val, "replace me")
	}

	if isLazy(meta.ID) {
		errors = append(errors, "metadata.id contains lazy value")
	}
	if isLazy(meta.Title) {
		errors = append(errors, "metadata.title contains lazy value")
	}
	if isLazy(meta.Summary) {
		errors = append(errors, "metadata.summary contains lazy value")
	}
	if len(meta.Summary) < 10 && !isLazy(meta.Summary) {
		errors = append(errors, "metadata.summary is too short (< 10 chars)")
	}
	if isLazy(meta.Type) {
		errors = append(errors, "metadata.type contains lazy value")
	}

	result := &domain.ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}

	return result
}

func parseCodeModel(data []byte) (*domain.CodeModel, error) {
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid format")
	}

	fontmatter := parts[1]
	body := parts[2]

	var meta domain.CodeModelMeta
	if err := yaml.Unmarshal([]byte(fontmatter), &meta); err != nil {
		return nil, err
	}

	return &domain.CodeModel{
		MetaData: meta,
		Body:     body,
	}, nil
}
