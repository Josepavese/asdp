package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Josepavese/asdp/engine/usecase"
)

type Server struct {
	queryUC     *usecase.QueryContextUseCase
	syncUC      *usecase.SyncModelUseCase
	scaffoldUC  *usecase.ScaffoldUseCase
	initAgentUC *usecase.InitAgentUseCase
}

func NewServer(queryUC *usecase.QueryContextUseCase, syncUC *usecase.SyncModelUseCase, scaffoldUC *usecase.ScaffoldUseCase, initAgentUC *usecase.InitAgentUseCase) *Server {
	return &Server{
		queryUC:     queryUC,
		syncUC:      syncUC,
		scaffoldUC:  scaffoldUC,
		initAgentUC: initAgentUC,
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
			Version: "0.1.6",
		},
	}, nil
}

func (s *Server) handleListTools() (*ListToolsResult, *RpcError) {
	return &ListToolsResult{
		Tools: []ToolDefinition{
			{
				Name:        "asdp_query_context",
				Description: "Retrieve the ASDP context (Spec, Model, Freshness) for a given path.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the module (default: .)",
						},
					},
					"required": []string{"path"},
				},
			},
			{
				Name:        "asdp_sync_codemodel",
				Description: "Automatically scans the source code and updates the codemodel.md file with fresh symbols and hash.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the module (default: .)",
						},
					},
					"required": []string{"path"},
				},
			},
			{
				Name:        "asdp_scaffold",
				Description: "Create a new ASDP-compliant module with codespec.md and codemodel.md.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type":        "string",
							"description": "The name of the module/directory.",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"description": "Module type (library, service, app). Default: library",
						},
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Parent directory for the new module. Default: .",
						},
					},
					"required": []string{"name"},
				},
			},
			{
				Name:        "asdp_init_agent",
				Description: "Copies ASDP Agent rules, skills, and workflows from global storage into the local project (.agent folder).",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Project root path. Default: .",
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
