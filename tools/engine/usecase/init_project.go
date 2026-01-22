package usecase

import (
	"fmt"
)

type InitProjectUseCase struct {
	fs       *InitAgentUseCase
	syncTree *SyncTreeUseCase
	scaffold *ScaffoldUseCase
}

func NewInitProjectUseCase(fs *InitAgentUseCase, syncTree *SyncTreeUseCase, scaffold *ScaffoldUseCase) *InitProjectUseCase {
	return &InitProjectUseCase{
		fs:       fs,
		syncTree: syncTree,
		scaffold: scaffold,
	}
}

func (uc *InitProjectUseCase) Execute(projectPath, codePath, title, summary, context string) (string, error) {
	if len(title) < 3 {
		return "", fmt.Errorf("title must be at least 3 characters long")
	}
	if len(summary) < 10 {
		return "", fmt.Errorf("summary must be at least 10 characters long")
	}
	if len(context) < 20 {
		return "", fmt.Errorf("context must be at least 20 characters long")
	}

	// 1. Initialize Agent at Project Root
	agentRes, err := uc.fs.Execute(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to init agent: %w", err)
	}

	// 2. Sync CodeTree at Code Root
	// If codePath is empty or ".", we use projectPath
	if codePath == "" || codePath == "." {
		codePath = projectPath
	}

	tree, err := uc.syncTree.Execute(codePath)
	if err != nil {
		return "", fmt.Errorf("failed to sync codetree: %w", err)
	}

	// 3. Ensure Code Root is a Module (Scaffold)
	// We scaffold with name="." in the codePath to ensure it has codespec/codemodel
	scaffoldRes, err := uc.scaffold.Execute(ScaffoldParams{
		Name:    ".",
		Path:    codePath,
		Type:    "module",
		Title:   title,
		Summary: summary,
		Context: context,
	})
	if err != nil {
		// Log but don't fail if scaffold fails (e.g. files already exist)
		scaffoldRes = fmt.Sprintf("Scaffold skipped or failed: %v", err)
	}

	return fmt.Sprintf("Project Anchor established.\n- Agent: %s\n- Tree: Root at %s (%d components)\n- Scaffold: %s",
		agentRes, codePath, len(tree.MetaData.Components), scaffoldRes), nil
}
