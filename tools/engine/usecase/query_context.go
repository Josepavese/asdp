package usecase

import (
	"fmt"
	"path/filepath"

	"github.com/Josepavese/asdp/engine/domain"
)

type QueryContextUseCase struct {
	fs     domain.FileSystem
	hasher domain.ContentHasher
}

func NewQueryContextUseCase(fs domain.FileSystem, hasher domain.ContentHasher) *QueryContextUseCase {
	return &QueryContextUseCase{fs: fs, hasher: hasher}
}

type ContextResponse struct {
	Path      string            `json:"path"`
	Summary   string            `json:"summary"`
	Freshness FreshnessStatus   `json:"freshness"`
	Spec      *domain.CodeSpec  `json:"spec,omitempty"`
	Model     *domain.CodeModel `json:"model,omitempty"`
}

type FreshnessStatus struct {
	Status      string `json:"status"` // "fresh", "stale", "unknown"
	Reason      string `json:"reason,omitempty"`
	CurrentHash string `json:"current_hash,omitempty"`
	DocHash     string `json:"doc_hash,omitempty"`
}

func (uc *QueryContextUseCase) Execute(path string) (*ContextResponse, error) {
	absPath, err := validateAndExpandPath(path)
	if err != nil {
		return nil, err
	}
	path = absPath

	resp := &ContextResponse{
		Path:      path,
		Freshness: FreshnessStatus{Status: "unknown"},
	}

	// 1. Read CodeSpec
	specBytes, err := uc.fs.ReadFile(filepath.Join(path, "codespec.md"))
	if err == nil {
		spec, err := parseCodeSpec(specBytes)
		if err == nil {
			resp.Spec = spec
			resp.Summary = spec.MetaData.Summary
			resp.Spec.Body = ""
		}
	}

	// 2. Read CodeModel
	modelBytes, err := uc.fs.ReadFile(filepath.Join(path, "codemodel.md"))
	if err == nil {
		model, err := parseCodeModel(modelBytes)
		if err == nil {
			resp.Model = model
			resp.Model.Body = ""

			// Optimize Payload: Strip docstrings to prevent JSON truncation
			for i := range resp.Model.MetaData.Symbols {
				resp.Model.MetaData.Symbols[i].Docstring = ""
			}
		}
	}

	// 3. Check Freshness
	if resp.Model != nil {
		realHash, err := uc.hasher.HashDir(path)
		if err == nil {
			docHash := resp.Model.MetaData.Integrity.SrcHash
			resp.Freshness.CurrentHash = realHash
			resp.Freshness.DocHash = docHash

			if docHash == "PENDING_FIRST_SYNC" {
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
