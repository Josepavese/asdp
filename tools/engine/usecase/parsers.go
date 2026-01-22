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

// Basic structural check only
func IsSpecParsable(data []byte) bool {
	parts := strings.SplitN(string(data), "---", 3)
	return len(parts) >= 3
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
