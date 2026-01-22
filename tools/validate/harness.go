package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Server represents a running ASDP MCP server process.
type Server struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
	stderr bytes.Buffer
}

// StartServer builds and starts the MCP server.
// returns the Server instance and a cleanup function.
func StartServer(t *testing.T, projectRoot string) (*Server, func()) {
	// 1. Build the server binary
	binPath := filepath.Join(os.TempDir(), "asdp-mcp-server-test")
	buildCmd := exec.Command("go", "build", "-o", binPath, "tools/mcp-server/cmd/asdp-mcp-server/main.go")
	buildCmd.Dir = projectRoot
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\n%s", err, out)
	}

	// 2. Start the process
	cmd := exec.Command(binPath)
	cmd.Dir = filepath.Join(projectRoot, "tests/sandbox") // Run inside a sandbox dir

	// Create sandbox if not exists
	os.MkdirAll(cmd.Dir, 0755)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("Failed to get stdin pipe: %v", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("Failed to get stdout pipe: %v", err)
	}

	// Capture stderr for debugging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server process: %v", err)
	}

	// Wait for startup message (optional, but good for stability)
	// Actually typical MCP server starts silent or logs to stderr?
	// Our main.go prints to Stderr: "ASDP MCP Server v... started."
	// We can check that later if needed, but for now we trust it starts.

	srv := &Server{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewScanner(stdoutPipe),
		stderr: stderr,
	}

	cleanup := func() {
		cmd.Process.Kill()
		os.Remove(binPath)
		// Optionally remove sandbox
		// os.RemoveAll(cmd.Dir)
	}

	return srv, cleanup
}

// CallTool sends a JSON-RPC request to execute a tool.
func (s *Server) CallTool(t *testing.T, name string, args map[string]interface{}) map[string]interface{} {
	id := time.Now().UnixNano()
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      name,
			"arguments": args,
		},
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	// Write request line-delimited
	_, err = fmt.Fprintln(s.stdin, string(reqBytes))
	if err != nil {
		t.Fatalf("Failed to write to stdin: %v", err)
	}

	// Read response
	// The server might output logs? But standard MCP should be line-delimited JSON on stdout.
	// Our server uses `json.NewEncoder(os.Stdout).Encode(resp)`?
	// Wait, `mcp.NewServer` handling.
	// `mcpServer.Serve()` reads stdin and writes stdout.

	if !s.stdout.Scan() {
		t.Fatalf("Server stdout closed unexpectedly. Stderr: %s", s.stderr.String())
	}
	respLine := s.stdout.Text()

	var resp map[string]interface{}
	if err := json.Unmarshal([]byte(respLine), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v\nLine: %s", err, respLine)
	}

	// Check for error field in JSON-RPC
	if errObj, ok := resp["error"]; ok {
		t.Fatalf("Tool execution returned error: %v", errObj)
	}

	result, ok := resp["result"].(map[string]interface{})
	if !ok {
		// It might be nested content?
		// MCP spec: result: { content: [...] }
		// Wait, if "error" is not present, "result" SHOULD be present.
		// If neither, invalid JSON-RPC.
		t.Fatalf("Response missing 'result' field: %s", respLine)
	}

	// Extract Content Text for convenience?
	// usually returning the whole result object is safer for specific assays.
	return result
}

// Assertions

func AssertFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist at %s, but it does not.", path)
	}
}

func AssertFileContent(t *testing.T, path string, contains string) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", path, err)
	}
	if !strings.Contains(string(content), contains) {
		t.Errorf("File %s does not contain expected string '%s'.\nContent: %s", path, contains, string(content))
	}
}

func AssertModTimeRecent(t *testing.T, path string) {
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Failed to stat file %s: %v", path, err)
	}
	if time.Since(info.ModTime()) > 5*time.Second {
		t.Errorf("File %s is stale (modtime %v is > 5s ago)", path, info.ModTime())
	}
}
