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
		AssertFileExists(t, filepath.Join(sandboxDir, ".agent/rules/managing-asdp-modules.asdp.md"))
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
		}
		srv.CallTool(t, "asdp_init_project", args)

		AssertFileExists(t, filepath.Join(projectDir, ".agent/rules/managing-asdp-modules.asdp.md"))
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
}
