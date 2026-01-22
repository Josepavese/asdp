package usecase

import (
	"fmt"

	"github.com/Josepavese/asdp/engine/domain"
)

type FunctionInfoResponse struct {
	Symbol  domain.Symbol          `json:"symbol"`
	Code    string                 `json:"code"`
	Context domain.ContextResponse `json:"context"`
}

type GetFunctionInfoUseCase struct {
	fs     domain.FileSystem
	parser domain.ASTParser
	hasher domain.ContentHasher
	config domain.Config
}

func NewGetFunctionInfoUseCase(fs domain.FileSystem, parser domain.ASTParser, hasher domain.ContentHasher, config domain.Config) *GetFunctionInfoUseCase {
	return &GetFunctionInfoUseCase{
		fs:     fs,
		parser: parser,
		hasher: hasher,
		config: config,
	}
}

func (uc *GetFunctionInfoUseCase) Execute(modulePath string, symbolName string) (*FunctionInfoResponse, error) {
	absPath, err := validateAndExpandPath(modulePath)
	if err != nil {
		return nil, err
	}
	modulePath = absPath

	// 1. Get Context (provides model and spec)
	// We could instantiate QueryContextUseCase here or just reuse logic.
	// Reusing logic is cleaner if we had a service, but let's keep it simple.
	queryUC := NewQueryContextUseCase(uc.fs, uc.hasher, uc.config)
	ctx, err := queryUC.Execute(modulePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get module context: %w", err)
	}

	// 2. Find symbol in model
	var targetSymbol *domain.Symbol
	for _, sym := range ctx.Model.MetaData.Symbols {
		if sym.Name == symbolName {
			targetSymbol = &sym
			break
		}
	}

	if targetSymbol == nil {
		return nil, fmt.Errorf("symbol %s not found in module %s", symbolName, modulePath)
	}

	// 3. Extract body
	body, err := uc.parser.GetSymbolBody(modulePath, *targetSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to extract symbol body: %w", err)
	}

	return &FunctionInfoResponse{
		Symbol:  *targetSymbol,
		Code:    body,
		Context: *ctx,
	}, nil
}
