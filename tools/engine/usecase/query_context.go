package usecase

import (
	"fmt"
	"path/filepath"

	"github.com/Josepavese/asdp/engine/domain"
)

type QueryContextUseCase struct {
	fs     domain.FileSystem
	hasher domain.ContentHasher
	config domain.Config
}

func NewQueryContextUseCase(fs domain.FileSystem, hasher domain.ContentHasher, config domain.Config) *QueryContextUseCase {
	return &QueryContextUseCase{fs: fs, hasher: hasher, config: config}
}

// ContextResponse moved to domain

func (uc *QueryContextUseCase) Execute(path string) (*domain.ContextResponse, error) {
	absPath, err := validateAndExpandPath(path)
	if err != nil {
		return nil, err
	}
	path = absPath

	resp := &domain.ContextResponse{
		Path:      path,
		Freshness: domain.Freshness{Status: "unknown"},
	}

	// 1. Read CodeSpec
	specBytes, err := uc.fs.ReadFile(filepath.Join(path, "codespec.md"))
	if err == nil {
		spec, err := parseCodeSpec(specBytes)
		if err == nil {
			resp.Spec = *spec
			resp.Summary = spec.MetaData.Summary
			resp.Spec.Body = ""

			// Policy-based validation is now centralized in ValidateProjectUseCase
			// We only do a basic structural check here if needed.
		}
	}

	// 2. Read CodeModel
	modelBytes, err := uc.fs.ReadFile(filepath.Join(path, "codemodel.md"))
	if err == nil {
		model, err := parseCodeModel(modelBytes)
		if err == nil {
			resp.Model = *model
			resp.Model.Body = ""

			// Optimize Payload: Strip docstrings to prevent JSON truncation
			for i := range resp.Model.MetaData.Symbols {
				resp.Model.MetaData.Symbols[i].Docstring = ""
			}
		}
	}

	// 3. Check Freshness
	if resp.Model.MetaData.ASDPVersion != "" {
		realHash, err := uc.hasher.HashDir(path)
		if err == nil {
			docHash := resp.Model.MetaData.Integrity.SrcHash
			resp.Freshness.CurrentHash = realHash
			resp.Freshness.DocHash = docHash

			if docHash == uc.config.Sync.Model.FirstSyncHash {
				resp.Freshness.Status = "stale"
				resp.Freshness.Reason = "New module, never synced"
			} else if docHash != realHash {
				resp.Freshness.Status = "stale"
				resp.Freshness.Reason = "Source code changed"
			} else {
				resp.Freshness.Status = "fresh"
			}
		} else {
			resp.Freshness.Reason = fmt.Sprintf("Hashing failed: %v", err)
		}
	} else {
		resp.Freshness.Reason = "No codemodel.md found"
	}

	return resp, nil
}
