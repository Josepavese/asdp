package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFunctionalSuite(t *testing.T) {
	// Setup
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(filepath.Dir(wd))

	srv, cleanup := StartServer(t, projectRoot)
	defer cleanup()

	sandboxDir := filepath.Join(projectRoot, "tools/validate/sandbox")
	os.RemoveAll(sandboxDir)
	os.MkdirAll(sandboxDir, 0755)
	defer os.RemoveAll(sandboxDir)

	// SCENARIO 1: SCAFFOLDING
	t.Run("Scaffolding", func(t *testing.T) {
		args := map[string]interface{}{
			"name":    "mymodule",
			"path":    sandboxDir,
			"title":   "My Module",
			"summary": "This is a test module",
			"context": "Context for testing",
			"type":    "library",
		}

		result := srv.CallTool(t, "asdp_scaffold", args)

		content := result["content"].([]interface{})
		text := content[0].(map[string]interface{})["text"].(string)
		if text == "" {
			t.Error("Scaffold returned empty output")
		}

		moduleDir := filepath.Join(sandboxDir, "mymodule")
		AssertFileExists(t, filepath.Join(moduleDir, "codespec.md"))
		AssertFileExists(t, filepath.Join(moduleDir, "codemodel.md"))

		AssertFileContent(t, filepath.Join(moduleDir, "codespec.md"), "title: \"My Module\"")
		AssertModTimeRecent(t, filepath.Join(moduleDir, "codespec.md"))
	})

	// SCENARIO 2: SYNC MODEL
	t.Run("Sync Model", func(t *testing.T) {
		moduleDir := filepath.Join(sandboxDir, "mymodule")
		srcFile := filepath.Join(moduleDir, "main.go")
		os.WriteFile(srcFile, []byte("package main\nfunc foo() {}"), 0644)

		args := map[string]interface{}{
			"path": moduleDir,
		}
		result := srv.CallTool(t, "asdp_sync_codemodel", args)

		content := result["content"].([]interface{})
		jsonStr := content[0].(map[string]interface{})["text"].(string)

		if !strings.Contains(jsonStr, "updated") && !strings.Contains(jsonStr, "refreshed_metadata") {
			t.Errorf("Expected updated status, got: %s", jsonStr)
		}

		AssertFileContent(t, filepath.Join(moduleDir, "codemodel.md"), "src_hash:")
		AssertModTimeRecent(t, filepath.Join(moduleDir, "codemodel.md"))
	})

	// SCENARIO 3: SYNC TREE
	t.Run("Sync Tree", func(t *testing.T) {
		args := map[string]interface{}{
			"path": sandboxDir,
		}
		srv.CallTool(t, "asdp_sync_codetree", args)

		AssertFileExists(t, filepath.Join(sandboxDir, "codetree.md"))
		AssertFileContent(t, filepath.Join(sandboxDir, "codetree.md"), "mymodule")
	})

	// SCENARIO 4: VALIDATE
	t.Run("Validate Errors", func(t *testing.T) {
		moduleDir := filepath.Join(sandboxDir, "mymodule")
		specPath := filepath.Join(moduleDir, "codespec.md")

		oldContent, _ := os.ReadFile(specPath)
		newContent := string(oldContent) + "\n TODO: Fix this"
		os.WriteFile(specPath, []byte(newContent), 0644)

		args := map[string]interface{}{
			"path": sandboxDir,
		}
		result := srv.CallTool(t, "asdp_validate", args)

		content := result["content"].([]interface{})
		jsonStr := content[0].(map[string]interface{})["text"].(string)

		if !strings.Contains(jsonStr, "\"is_valid\": false") {
			t.Errorf("Expected validation failure, got: %s", jsonStr)
		}
	})

	// SCENARIO 5: QUERY CONTEXT
	t.Run("Query Context", func(t *testing.T) {
		args := map[string]interface{}{
			"path": sandboxDir + "/mymodule",
		}
		result := srv.CallTool(t, "asdp_query_context", args)

		content := result["content"].([]interface{})
		jsonStr := content[0].(map[string]interface{})["text"].(string)

		if !strings.Contains(jsonStr, "My Module") {
			t.Errorf("Expected query result to contain title, got: %s", jsonStr)
		}
	})

	// SCENARIO 6: INIT AGENT
	t.Run("Init Agent", func(t *testing.T) {
		args := map[string]interface{}{
			"path": sandboxDir,
		}
		srv.CallTool(t, "asdp_init_agent", args)

		// Check for one known file from global assets
		AssertFileExists(t, filepath.Join(sandboxDir, ".agent/rules/managing-asdp-modules-asdp.md"))
	})

	// SCENARIO 7: INIT PROJECT
	t.Run("Init Project", func(t *testing.T) {
		projectDir := filepath.Join(sandboxDir, "newproject")
		os.MkdirAll(projectDir, 0755)

		// Note: since InitProject calls Scaffold Internally,
		// and Scaffold now requires Title/Summary/Context if strict,
		// we might need to update asdp_init_project to accept those or use defaults.
		// For now, InitProject skips scaffold if it fails.

		args := map[string]interface{}{
			"path":      projectDir,
			"code_path": projectDir,
			"title":     "Root Module Title for Testing",
			"summary":   "Root summary must be long enough for testing",
			"context":   "Root context must be very descriptive and long enough to pass the validation check.",
		}
		srv.CallTool(t, "asdp_init_project", args)

		AssertFileExists(t, filepath.Join(projectDir, ".agent/rules/managing-asdp-modules-asdp.md"))
		AssertFileExists(t, filepath.Join(projectDir, "codetree.md"))
	})

	// SCENARIO 8: FUNCTION INFO
	t.Run("Function Info", func(t *testing.T) {
		moduleDir := filepath.Join(sandboxDir, "mymodule")
		args := map[string]interface{}{
			"path":   moduleDir,
			"symbol": "foo",
		}
		result := srv.CallTool(t, "asdp_function_info", args)

		content := result["content"].([]interface{})
		jsonStr := content[0].(map[string]interface{})["text"].(string)

		if !strings.Contains(jsonStr, "func foo() {}") {
			t.Errorf("Expected function code, got: %s", jsonStr)
		}
		if !strings.Contains(jsonStr, "\"kind\": \"function\"") {
			t.Errorf("Expected symbol metadata, got: %s", jsonStr)
		}
	})
	// SCENARIO 9: EXCLUSIONS & VALIDATION (Deep Test)
	t.Run("Exclusions End-to-End", func(t *testing.T) {
		exclusionDir := filepath.Join(sandboxDir, "exclusion_e2e")
		os.MkdirAll(exclusionDir, 0755)

		// 0. Init Project (Creates codetree.md and codespec.md)
		srv.CallTool(t, "asdp_sync_codetree", map[string]interface{}{"path": exclusionDir})
		os.WriteFile(filepath.Join(exclusionDir, "codespec.md"), []byte("---\ntitle: Root\nsummary: Root context\n---\n## Context\nTest ctx"), 0644)

		// 1. Create a "Broken" Module (Significant but missing spec)
		// This should trigger a validation error.
		brokenDir := filepath.Join(exclusionDir, "broken_lib")
		os.MkdirAll(brokenDir, 0755)
		os.WriteFile(filepath.Join(brokenDir, "logic.go"), []byte("package broken"), 0644)

		// 2. Validate -> SHOULD FAIL
		resFail := srv.CallTool(t, "asdp_validate", map[string]interface{}{"path": exclusionDir})
		contentFail := resFail["content"].([]interface{})[0].(map[string]interface{})["text"].(string)

		if !strings.Contains(contentFail, "Missing required file: codespec.md") {
			t.Errorf("Expected validation error for broken_lib, got: %s", contentFail)
		}
		if !strings.Contains(contentFail, "asdp_manage_exclusions") {
			t.Errorf("Expected validation error to suggest exclusions, got: %s", contentFail)
		}

		// 3. Exclude 'broken_lib'
		srv.CallTool(t, "asdp_manage_exclusions", map[string]interface{}{
			"path":   exclusionDir,
			"target": "broken_lib",
			"action": "add",
		})

		// 4. Validate -> SHOULD PASS
		// The validator should now skip 'broken_lib' because it's excluded in codetree.md
		// Note: We need to ensure Validate logic actually READS exclusions.
		// Currently ProjectValidator.Execute checks 'shouldIgnoreDir', but that checked names starting with '.'.
		// It mentions exclusions in the error message, but DOES IT actually respect codetree exclusions during the walk?
		// I need to check `project_validator.go` implementation again.
		// If check pass, we are good. If not, I found a bug to fix!

		resPass := srv.CallTool(t, "asdp_validate", map[string]interface{}{"path": exclusionDir})
		contentPass := resPass["content"].([]interface{})[0].(map[string]interface{})["text"].(string)

		if strings.Contains(contentPass, "Missing required file: codespec.md") {
			t.Errorf("Validation still failed after exclusion! Output: %s", contentPass)
		}
		if !strings.Contains(contentPass, "\"is_valid\": true") {
			t.Errorf("Expected is_valid: true, got: %s", contentPass)
		}

		// 5. Test Branch Exclusion
		branchDir := filepath.Join(exclusionDir, "legacy/v1/nested")
		os.MkdirAll(branchDir, 0755)
		os.WriteFile(filepath.Join(branchDir, "old.go"), []byte("package old"), 0644)

		// Verify it fails first
		resBranchFail := srv.CallTool(t, "asdp_validate", map[string]interface{}{"path": exclusionDir})
		if !strings.Contains(resBranchFail["content"].([]interface{})[0].(map[string]interface{})["text"].(string), "Missing required file") {
			t.Error("Expected validation error for legacy branch")
		}

		// Exclude parent 'legacy'
		srv.CallTool(t, "asdp_manage_exclusions", map[string]interface{}{
			"path":   exclusionDir,
			"target": "legacy",
			"action": "add",
		})

		// Verify it passes
		resBranchPass := srv.CallTool(t, "asdp_validate", map[string]interface{}{"path": exclusionDir})
		jsonPass := resBranchPass["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
		if !strings.Contains(jsonPass, "\"is_valid\": true") {
			t.Errorf("Validation failed after branch exclusion. Output: %s", jsonPass)
		}
	})
}
