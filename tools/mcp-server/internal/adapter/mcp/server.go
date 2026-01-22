package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Josepavese/asdp/engine/domain"
	"github.com/Josepavese/asdp/engine/usecase"
	"github.com/Josepavese/asdp/validate/check"
)

type Server struct {
	queryUC       *usecase.QueryContextUseCase
	syncUC        *usecase.SyncModelUseCase
	scaffoldUC    *usecase.ScaffoldUseCase
	initAgentUC   *usecase.InitAgentUseCase
	syncTreeUC    *usecase.SyncTreeUseCase
	initProjectUC *usecase.InitProjectUseCase
	validateUC    *check.ValidateProjectUseCase
	functionUC    *usecase.GetFunctionInfoUseCase
	config        domain.Config
}

func NewServer(queryUC *usecase.QueryContextUseCase, syncUC *usecase.SyncModelUseCase, scaffoldUC *usecase.ScaffoldUseCase, initAgentUC *usecase.InitAgentUseCase, syncTreeUC *usecase.SyncTreeUseCase, initProjectUC *usecase.InitProjectUseCase, validateUC *check.ValidateProjectUseCase, functionUC *usecase.GetFunctionInfoUseCase, config domain.Config) *Server {
	return &Server{
		queryUC:       queryUC,
		syncUC:        syncUC,
		scaffoldUC:    scaffoldUC,
		initAgentUC:   initAgentUC,
		syncTreeUC:    syncTreeUC,
		initProjectUC: initProjectUC,
		validateUC:    validateUC,
		functionUC:    functionUC,
		config:        config,
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
				Description: "Retrieve the ASDP context (Spec, Model, Freshness) for a given absolute path. Result: Returns a JSON object containing the merged CodeSpec (intent), CodeModel (structure), and current freshness status, allowing an agent to quickly understand a module's contract and implementation.",
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
				Description: "Automatically scans the source code and updates the codemodel.md file. Result: Returns a SyncResult JSON with the count of symbols identified (functions, structs, interfaces including start/end lines) and the integrity hash of the source files.",
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
				Name:        "asdp_sync_codetree",
				Description: "Automatically scans the project directory and updates the codetree.md file. Result: Returns a JSON representation of the project hierarchy, listing all modules and their ASDP compliance status (presence of codespec/codemodel).",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE path to the project root or sub-directory.",
						},
					},
					"required": []string{"path"},
				},
			},
			{
				Name:        "asdp_scaffold",
				Description: "Create a new ASDP-compliant module or backfill missing files (codespec/codemodel) in an existing one. Safe to run on existing directories; will not overwrite existing files. Result: Returns a success message.",
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
						"title": map[string]interface{}{
							"type":        "string",
							"description": "Title of the module (e.g. 'User Authentication Service'). Required.",
						},
						"summary": map[string]interface{}{
							"type":        "string",
							"description": "Brief summary of the module's purpose. Required.",
						},
						"context": map[string]interface{}{
							"type":        "string",
							"description": "Detailed context and reasoning for this module. Required.",
						},
					},
					"required": []string{"name", "path", "title", "summary", "context"},
				},
			},
			{
				Name:        "asdp_init_agent",
				Description: "Copies ASDP Agent assets into the local project. Result: Returns a list of files copied into the .agent/ directory (Rules, Workflows, etc.) to enable ASDP-native agent behavior.",
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
			{
				Name:        "asdp_init_project",
				Description: "Unified initialization of an ASDP project. Result: Sets up .agent/ at project path AND anchors the CodeTree at a specified code path, ensuring the AI starts from the real code.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE project repository root path.",
						},
						"code_path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE path to where the actual code starts (e.g. /repo/tools).",
						},
					},
					"required": []string{"path", "code_path"},
				},
			},
			{
				Name:        "asdp_validate",
				Description: "Audit the ASDP project state. Returns a report of Errors (invalid state, integration blocking) and Warnings (staleness). Checks for mandatory files, strict content compliance, and synchronization freshness.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "ABSOLUTE path to the project root.",
						},
					},
					"required": []string{"path"},
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

	switch callParams.Name {
	case "asdp_query_context":
		path, _ := callParams.Arguments["path"].(string)
		ctx, err := s.queryUC.Execute(path)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		jsonBytes, _ := json.MarshalIndent(ctx, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
			IsError: ctx.Validation != nil && !ctx.Validation.IsValid,
		}, nil

	case "asdp_sync_codemodel":
		path, _ := callParams.Arguments["path"].(string)
		res, err := s.syncUC.Execute(path)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		jsonBytes, _ := json.MarshalIndent(res, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
		}, nil

	case "asdp_sync_codetree":
		path, _ := callParams.Arguments["path"].(string)
		res, err := s.syncTreeUC.Execute(path)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}

		// Recursive check for validation errors anywhere in the tree
		hasAnyError := false
		var checkErrors func(c []domain.Component)
		checkErrors = func(comps []domain.Component) {
			for _, c := range comps {
				if !c.IsValid {
					hasAnyError = true
					return
				}
				checkErrors(c.Children)
			}
		}
		if !res.MetaData.Root || res.MetaData.Components == nil {
			// This part is tricky because the root itself might have a spec (not yet in TreeMeta)
			// But SyncTree builds children.
		}
		checkErrors(res.MetaData.Components)

		jsonBytes, _ := json.MarshalIndent(res, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
			IsError: hasAnyError,
		}, nil

	case "asdp_scaffold":
		name, _ := callParams.Arguments["name"].(string)
		modType, _ := callParams.Arguments["type"].(string)
		path, _ := callParams.Arguments["path"].(string)
		title, _ := callParams.Arguments["title"].(string)
		summary, _ := callParams.Arguments["summary"].(string)
		context, _ := callParams.Arguments["context"].(string)
		if name == "" {
			return nil, &RpcError{Code: -32602, Message: "Missing required argument: name"}
		}
		if modType == "" {
			modType = "library"
		}
		resultMsg, err := s.scaffoldUC.Execute(usecase.ScaffoldParams{
			Name:    name,
			Type:    modType,
			Path:    path,
			Title:   title,
			Summary: summary,
			Context: context,
		})
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		return &CallToolResult{Content: []ToolContent{{Type: "text", Text: resultMsg}}}, nil

	case "asdp_init_agent":
		path, _ := callParams.Arguments["path"].(string)
		resultMsg, err := s.initAgentUC.Execute(path)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		return &CallToolResult{Content: []ToolContent{{Type: "text", Text: resultMsg}}}, nil

	case "asdp_init_project":
		path, _ := callParams.Arguments["path"].(string)
		codePath, _ := callParams.Arguments["code_path"].(string)
		title, _ := callParams.Arguments["title"].(string)
		summary, _ := callParams.Arguments["summary"].(string)
		context, _ := callParams.Arguments["context"].(string)

		resultMsg, err := s.initProjectUC.Execute(path, codePath, title, summary, context)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		return &CallToolResult{Content: []ToolContent{{Type: "text", Text: resultMsg}}}, nil

	case "asdp_validate":
		path, _ := callParams.Arguments["path"].(string)
		report, err := s.validateUC.Execute(path)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		jsonBytes, _ := json.MarshalIndent(report, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
			IsError: !report.IsValid,
		}, nil

	case "asdp_function_info":
		path, _ := callParams.Arguments["path"].(string)
		symbol, _ := callParams.Arguments["symbol"].(string)
		res, err := s.functionUC.Execute(path, symbol)
		if err != nil {
			return nil, &RpcError{Code: -32000, Message: err.Error()}
		}
		jsonBytes, _ := json.MarshalIndent(res, "", "  ")
		return &CallToolResult{
			Content: []ToolContent{{Type: "text", Text: string(jsonBytes)}},
		}, nil

	default:
		return nil, &RpcError{Code: -32601, Message: fmt.Sprintf("Tool not found: %s", callParams.Name)}
	}
}
