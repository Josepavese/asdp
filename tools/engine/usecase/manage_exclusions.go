package usecase

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Josepavese/asdp/engine/domain"
	"gopkg.in/yaml.v3"
)

type ManageExclusionsUseCase struct {
	fs         domain.FileSystem
	syncTreeUC *SyncTreeUseCase
}

func NewManageExclusionsUseCase(fs domain.FileSystem, syncTreeUC *SyncTreeUseCase) *ManageExclusionsUseCase {
	return &ManageExclusionsUseCase{
		fs:         fs,
		syncTreeUC: syncTreeUC,
	}
}

func (uc *ManageExclusionsUseCase) Execute(projectPath string, target string, action string) error {
	absPath, err := validateAndExpandPath(projectPath)
	if err != nil {
		return err
	}
	projectPath = absPath

	treePath := filepath.Join(projectPath, "codetree.md")
	data, err := uc.fs.ReadFile(treePath)
	if err != nil {
		return fmt.Errorf("codetree.md not found at project root. Please init project first.")
	}

	// Split frontmatter manually to preserve body
	parts := strings.SplitN(string(data), "---", 3)
	if len(parts) < 3 {
		return fmt.Errorf("codetree.md has invalid format (missing frontmatter delimiters)")
	}

	frontmatter := parts[1]
	body := parts[2]

	var meta domain.CodeTreeMeta
	if err := yaml.Unmarshal([]byte(frontmatter), &meta); err != nil {
		return fmt.Errorf("failed to parse codetree.md frontmatter: %w", err)
	}

	// Update Excludes
	newExcludes := []string{}
	seen := make(map[string]bool)

	// Populate initial list avoiding duplicates if malformed
	if meta.Excludes != nil {
		for _, e := range meta.Excludes {
			if !seen[e] {
				newExcludes = append(newExcludes, e)
				seen[e] = true
			}
		}
	}

	if action == "add" {
		if !seen[target] {
			newExcludes = append(newExcludes, target)
		}
	} else if action == "remove" {
		// Filter out
		filtered := []string{}
		for _, e := range newExcludes {
			if e != target {
				filtered = append(filtered, e)
			}
		}
		newExcludes = filtered
	} else {
		return fmt.Errorf("unknown action '%s': must be 'add' or 'remove'", action)
	}

	meta.Excludes = newExcludes

	// Marshal back
	fmBytes, err := yaml.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to list exclusions: %w", err)
	}

	// Write back file
	newContent := fmt.Sprintf("---\n%s---\n%s", string(fmBytes), body)
	if err := uc.fs.WriteFile(treePath, []byte(newContent)); err != nil {
		return fmt.Errorf("failed to write updated codetree.md: %w", err)
	}

	// Trigger SyncTree to refresh the view immediately with new exclusions
	_, err = uc.syncTreeUC.Execute(projectPath)
	if err != nil {
		return fmt.Errorf("exclusions saved but failed to refresh tree: %w", err)
	}

	return nil
}
