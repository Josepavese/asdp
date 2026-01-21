package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Josepavese/asdp/engine/domain"
	"github.com/Josepavese/asdp/engine/usecase"
)

type Server struct {
	queryUC     *usecase.QueryContextUseCase
	syncUC      *usecase.SyncModelUseCase
	scaffoldUC  *usecase.ScaffoldUseCase
	initAgentUC *usecase.InitAgentUseCase
	syncTreeUC  *usecase.SyncTreeUseCase
}

func NewServer(queryUC *usecase.QueryContextUseCase, syncUC *usecase.SyncModelUseCase, scaffoldUC *usecase.ScaffoldUseCase, initAgentUC *usecase.InitAgentUseCase, syncTreeUC *usecase.SyncTreeUseCase) *Server {
	return &Server{
		queryUC:     queryUC,
		syncUC:      syncUC,
		scaffoldUC:  scaffoldUC,
		initAgentUC: initAgentUC,
		syncTreeUC:  syncTreeUC,
	}
}

// Serve starts the JSON-RPC loop on Stdin/Stdout
func (s *Server) Serve() {
	scanner := bufio.NewScanner(os.Stdin)
	// Increase buffer size for large JSON payloads if necessary
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		s.handleMessage(line)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scanner error: %v", err)
	}
}

func (s *Server) handleMessage(data []byte) {
	var req JsonRpcRequest
	if err := json.Unmarshal(data, &req); err != nil {
		// Ignore parse errors (might be logs mixed in stdio)
		return
	}

	var resp interface{}
	var err *RpcError

	switch req.Method {
	case "initialize":
		resp, err = s.handleInitialize(req.Params)
	case "tools/list":
		resp, err = s.handleListTools()
	case "tools/call":
		resp, err = s.handleCallTool(req.Params)
	default:
		// Unknown method, but for notifications we shouldn't error.
		// For now, we only care about request/response cycle.
		if req.ID != nil {
			err = &RpcError{Code: -32601, Message: "Method not found"}
		}
	}

	if req.ID != nil {
		s.sendResponse(req.ID, resp, err)
	}
}

func (s *Server) sendResponse(id interface{}, result interface{}, rpcErr *RpcError) {
	resp := JsonRpcResponse{
		JSONRPC: "2.0",
		ID:      id,
	}
	if rpcErr != nil {
		resp.Error = rpcErr
	} else {
		resp.Result = result
	}

	bytes, _ := json.Marshal(resp)
	fmt.Printf("%s\n", string(bytes))
}

// --- Handlers ---

func (s *Server) handleInitialize(params json.RawMessage) (*InitializeResult, *RpcError) {
	return &InitializeResult{
		ProtocolVersion: "2024-11-05", // MCP Version
		Capabilities: struct {
			Tools *struct {
				ListChanged bool `json:"listChanged"`
			} `json:"tools,omitempty"`
		}{
			Tools: &struct {
				ListChanged bool `json:"listChanged"`
			}{ListChanged: false},
		},
		ServerInfo: struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    "asdp-mcp-server",
			Version: domain.Version,
		},
	}, nil
}

func (s *Server) handleListTools() (*ListToolsResult, *RpcError) {
	return &ListToolsResult{
		Tools: []ToolDefinition{
			{
				Name:        "asdp_query_context",
				Description: "Retrieve the ASDP context (Spec, Model, Freshness) for a given absolute path.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE path to the module (e.g. /home/user/project/module)",
						},
					},
					"required": []string{"path"},
				},
			},
			{
				Name:        "asdp_sync_codemodel",
				Description: "Automatically scans the source code and updates the codemodel.md file. Requires an absolute path.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE path to the module (e.g. /home/user/project/module)",
						},
					},
					"required": []string{"path"},
				},
			},
			{
				Name:        "asdp_scaffold",
				Description: "Create a new ASDP-compliant module. Use name='.' to initialize the current directory (in-place).",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "Module name. Use '.' to scaffold directly in the provided path.",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Module type (library, service, app). Default: library",
						},
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE parent directory (or target directory if name='.').",
						},
					},
					"required": []string{"name", "path"},
				},
			},
			{
				Name:        "asdp_init_agent",
				Description: "Copies ASDP Agent assets into the local project. Specify the project root as an absolute path.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE project root path (e.g. /home/user/project)",
						},
					},
				},
			},
		},
	}, nil
}

func (s *Server) handleCallTool(params json.RawMessage) (*CallToolResult, *RpcError) {
	var callParams CallToolParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, &RpcError{Code: -32700, Message: "Parse error"}
	}

	if callParams.Name == "asdp_query_context" {
		path, ok := callParams.Arguments["path"].(string)
		if !ok || path == "" {
			path = "."
		}

		ctx, err := s.queryUC.Execute(path)
		if err != nil {
			return &CallToolResult{
				IsError: true,
				Content: []ToolContent{{Type: "text", Text: err.Error()}},
			}, nil
		}

		jsonBytes, _ := json.MarshalIndent(ctx, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
		}, nil
	}

	if callParams.Name == "asdp_sync_codemodel" {
		path, ok := callParams.Arguments["path"].(string)
		if !ok || path == "" {
			path = "."
		}

		res, err := s.syncUC.Execute(path)
		if err != nil {
			return &CallToolResult{
				IsError: true,
				Content: []ToolContent{{Type: "text", Text: err.Error()}},
			}, nil
		}

		jsonBytes, _ := json.MarshalIndent(res, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
		}, nil
	}

	if callParams.Name == "asdp_scaffold" {
		name, _ := callParams.Arguments["name"].(string)
		modType, _ := callParams.Arguments["type"].(string)
		path, _ := callParams.Arguments["path"].(string)

		if name == "" {
			return nil, &RpcError{Code: -32602, Message: "Missing required argument: name"}
		}
		if modType == "" {
			modType = "library"
		}

		params := usecase.ScaffoldParams{
			Name: name,
			Type: modType,
			Path: path,
		}

		resultMsg, err := s.scaffoldUC.Execute(params)
		if err != nil {
			return &CallToolResult{
				IsError: true,
				Content: []ToolContent{{Type: "text", Text: err.Error()}},
			}, nil
		}

		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: resultMsg}},
		}, nil
	}

	if callParams.Name == "asdp_init_agent" {
		path, ok := callParams.Arguments["path"].(string)
		if !ok || path == "" {
			path = "."
		}

		resultMsg, err := s.initAgentUC.Execute(path)
		if err != nil {
			return &CallToolResult{
				IsError: true,
				Content: []ToolContent{{Type: "text", Text: err.Error()}},
			}, nil
		}

		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: resultMsg}},
		}, nil
	}

	return nil, &RpcError{Code: -32601, Message: fmt.Sprintf("Tool not found: %s", callParams.Name)}
}
