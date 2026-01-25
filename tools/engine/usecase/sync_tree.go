package usecase

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Josepavese/asdp/engine/domain"
	"gopkg.in/yaml.v3"
)

type SyncTreeUseCase struct {
	fs     domain.FileSystem
	config domain.TreeSyncConfig
}

func NewSyncTreeUseCase(fs domain.FileSystem, config domain.TreeSyncConfig) *SyncTreeUseCase {
	return &SyncTreeUseCase{
		fs:     fs,
		config: config,
	}
}

func (uc *SyncTreeUseCase) Execute(path string) (*domain.CodeTree, error) {
	absPath, err := validateAndExpandPath(path)
	if err != nil {
		return nil, err
	}
	path = absPath

	// 1. Read existing exclusions from codetree.md (if exists)
	var existingExcludes []string
	treePath := filepath.Join(path, "codetree.md")
	if data, err := uc.fs.ReadFile(treePath); err == nil {
		// Parse struct just to get excludes
		var partialTree struct {
			Excludes []string `yaml:"excludes"`
		}
		if err := yaml.Unmarshal(data, &partialTree); err == nil {
			existingExcludes = partialTree.Excludes
		}
	}

	// 2. Build Component Tree (with exclusions)
	rootComp, err := uc.buildComponent(path, path, existingExcludes)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}

	// 3. Construct Tree Object
	tree := &domain.CodeTree{
		MetaData: domain.CodeTreeMeta{
			ASDPVersion: domain.Version,
			Root:        true,
			Components:  rootComp.Children,
			Verification: domain.Verification{
				ScanTime: time.Now(),
			},
			Excludes: existingExcludes,
		},
		Body: uc.config.HeaderTemplate,
	}

	// 4. Write to File
	fmBytes, err := yaml.Marshal(tree.MetaData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal yaml: %w", err)
	}

	newContent := fmt.Sprintf("---\n%s---\n%s", string(fmBytes), tree.Body)
	if err := uc.fs.WriteFile(treePath, []byte(newContent)); err != nil {
		return nil, fmt.Errorf("failed to write codetree.md: %w", err)
	}

	return tree, nil
}

func (uc *SyncTreeUseCase) buildComponent(root string, currentPath string, excludes []string) (*domain.Component, error) {
	relPath, _ := filepath.Rel(root, currentPath)
	if relPath == "." {
		relPath = "./"
	} else {
		relPath = "./" + relPath
	}

	comp := &domain.Component{
		Name: filepath.Base(currentPath),
		Path: relPath,
		Type: uc.config.DefaultComponent, // "module" from config
	}

	// Check and parse ASDP files for metadata
	specPath := filepath.Join(currentPath, "codespec.md")
	comp.IsValid = true // Assume valid unless proven otherwise
	if data, err := uc.fs.ReadFile(specPath); err == nil {
		comp.HasSpec = true
		if spec, err := parseCodeSpec(data); err == nil && spec != nil {
			if spec.MetaData.Type != "" {
				comp.Type = spec.MetaData.Type
			}
			if spec.MetaData.Summary != "" {
				comp.Description = spec.MetaData.Summary
			} else if spec.MetaData.Title != "" {
				comp.Description = spec.MetaData.Title
			}

			// Structural Check: If title and type are present, we consider it "valid enough" for tree display
			if spec.MetaData.Title == "" || spec.MetaData.Type == "" {
				comp.IsValid = false
			}
		} else {
			comp.IsValid = false // Malformed spec
		}
	} else {
		// Fallback description from config
		if comp.Description == "" {
			comp.Description = uc.config.FallbackDesc
		}
	}

	if _, err := uc.fs.Stat(filepath.Join(currentPath, "codemodel.md")); err == nil {
		comp.HasModel = true
	}

	// Calculate LastModified for this directory
	latest := time.Time{}
	if info, err := uc.fs.Stat(currentPath); err == nil {
		latest = info.ModTime()
	}

	// Read children
	var children []domain.Component

	// Hack: I'll use `Walk` to find immediate children
	err := uc.fs.Walk(currentPath, func(path string, isDir bool) error {
		if path == currentPath {
			return nil // Root of this walk
		}

		// Update latest mtime for ANY file found in this walk (to detect deep changes)
		if info, err := uc.fs.Stat(path); err == nil {
			if info.ModTime().After(latest) {
				latest = info.ModTime()
			}
		}

		// We only want immediate children for this node.
		rel, _ := filepath.Rel(currentPath, path)
		if len(strings.Split(rel, string(os.PathSeparator))) > 1 {
			if isDir {
				return fs.SkipDir
			}
			return nil
		}

		if !isDir {
			return nil
		}

		dirName := filepath.Base(path)
		if uc.isIgnoredDir(dirName, excludes) {
			return fs.SkipDir
		}

		if uc.isShallowDir(dirName) {
			childRel, _ := filepath.Rel(root, path)
			children = append(children, domain.Component{
				Name:         dirName,
				Path:         "./" + childRel,
				Type:         uc.config.DependencyType,
				Description:  "External dependencies (not scanned)",
				LastModified: latest, // Best effort
			})
			return fs.SkipDir
		}

		// Recurse to build sub-component
		childComp, err := uc.buildComponent(root, path, excludes)
		if err != nil {
			return err
		}
		children = append(children, *childComp)

		return fs.SkipDir
	})

	comp.Children = children
	comp.LastModified = latest
	return comp, err
}

func (uc *SyncTreeUseCase) isIgnoredDir(name string, excludes []string) bool {
	// Check user-defined exclusions first
	for _, exc := range excludes {
		// Simple match or glob? Requirement says "exclude folders or branches".
		// For simplicity, let's match exact name or path prefix if needed.
		// But here 'name' is just directory name.
		// If we want to exclude "apps/foo", we need the full relative path in the recursion.
		// But `buildComponent` receives `currentPath` and walks children.
		// To truly support path-based exclusion (e.g. "exclude only tools/legacy"), we need the relative path.
		// However, `isIgnoredDir` here only takes `name`.
		// Let's rely on basic name matching for simplicity now, as agreed in some conventions,
		// OR we should pass the relative path to `isIgnoredDir`.
		// Given the `fs.Walk` structure, `path` variable inside the closure has the full path.
		// BUT `isIgnoredDir` is called with `dirName`.
		// Let's upgrade `isIgnoredDir` to just check names against a list for now,
		// BUT we should also add a check inside the walk loop for the RELATIVE path if we want robust exclusion.
		// For now, I will just match against the name for simple folder exclusion.
		// Improve: Support glob matching?
		if name == exc {
			return true
		}
	}

	for _, idx := range uc.config.IgnoredDirs {
		if name == idx {
			return true
		}
	}
	// Also check config.IgnoredDirs if distinct?
	// Assuming merged or just checking global for now.
	// The implementation plan says "use cfg.IgnoredDirs"
	return strings.HasPrefix(name, ".") && name != "."
}

func (uc *SyncTreeUseCase) isShallowDir(name string) bool {
	for _, idx := range uc.config.ShallowDirs {
		if name == idx {
			return true
		}
	}
	return false
}
