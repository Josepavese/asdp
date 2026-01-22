package domain

import "time"

// --- CodeSpec (Intent) ---

type CodeSpec struct {
	MetaData CodeSpecMeta `yaml:",inline"`
	Body     string       `yaml:"-"` // Markdown content
}

type CodeSpecMeta struct {
	ASDPVersion  string        `yaml:"asdp_version"`
	LastModified time.Time     `yaml:"last_modified"`
	ID           string        `yaml:"id"`
	Type         string        `yaml:"type"` // library, app, service
	Title        string        `yaml:"title"`
	Summary      string        `yaml:"summary"`
	Capabilities []string      `yaml:"capabilities"`
	Dependencies []Dependency  `yaml:"dependencies"`
	Requirements []Requirement `yaml:"requirements"`
	Exports      []string      `yaml:"exports"`
}

type Dependency struct {
	Module string `yaml:"module"`
	Reason string `yaml:"reason"`
}

type Requirement struct {
	ID       string `yaml:"id"`
	Desc     string `yaml:"desc"`
	Priority string `yaml:"priority"`
}

// --- CodeModel (Structure) ---

type CodeModel struct {
	MetaData CodeModelMeta `yaml:",inline"`
	Body     string        `yaml:"-"`
}

type CodeModelMeta struct {
	ASDPVersion string    `yaml:"asdp_version"`
	Integrity   Integrity `yaml:"integrity"`
	Symbols     []Symbol  `yaml:"symbols"`
}

type Integrity struct {
	SrcHash      string    `yaml:"src_hash"`
	Algorithm    string    `yaml:"algorithm"`
	LastModified time.Time `yaml:"last_modified"`
	CheckedAt    time.Time `yaml:"checked_at"`
}

type Symbol struct {
	Name      string `yaml:"name" json:"name"`
	Kind      string `yaml:"kind" json:"kind"` // function, struct, class
	Exported  bool   `yaml:"exported" json:"exported"`
	Line      int    `yaml:"line" json:"line"`
	LineEnd   int    `yaml:"line_end" json:"line_end"`
	FilePath  string `yaml:"file_path,omitempty" json:"file_path,omitempty"`
	Signature string `yaml:"signature" json:"signature"`
	Docstring string `yaml:"docstring,omitempty" json:"docstring,omitempty"`
	Parent    string `yaml:"parent,omitempty" json:"parent,omitempty"`
}

// --- CodeTree (Hierarchy) ---

type CodeTree struct {
	MetaData CodeTreeMeta `yaml:",inline"`
	Body     string       `yaml:"-"`
}

type CodeTreeMeta struct {
	ASDPVersion  string       `yaml:"asdp_version"`
	Root         bool         `yaml:"root"`
	Components   []Component  `yaml:"components"`
	Verification Verification `yaml:"verification"`
}

type Component struct {
	Name         string      `yaml:"name"`
	Type         string      `yaml:"type"`
	Path         string      `yaml:"path"`
	Description  string      `yaml:"description"`
	LastModified time.Time   `yaml:"last_modified"`
	HasSpec      bool        `yaml:"has_spec"`
	HasModel     bool        `yaml:"has_model"`
	IsValid      bool        `yaml:"is_valid"`
	Children     []Component `yaml:"children,omitempty"`
}

type Verification struct {
	ScanTime time.Time `yaml:"scan_time"`
}

// --- Validation (Quality) ---

type ValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
}

type Freshness struct {
	Status      string `json:"status"` // "fresh", "stale", "unknown"
	Reason      string `json:"reason,omitempty"`
	CurrentHash string `json:"current_hash,omitempty"`
	DocHash     string `json:"doc_hash,omitempty"`
}

// ContextResponse is the DTO for QueryContext
type ContextResponse struct {
	Path       string            `json:"path"`
	Summary    string            `json:"summary"`
	Freshness  Freshness         `json:"freshness"`
	Validation *ValidationResult `json:"validation,omitempty"`
	Spec       CodeSpec          `json:"spec"`
	Model      CodeModel         `json:"model"`
}
