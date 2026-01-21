package domain

import "time"

// FileSystem abstraction for testability
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	MkdirAll(path string) error
	Walk(root string, fn func(path string, isDir bool) error) error
	Stat(path string) (FileInfo, error)
}

type FileInfo interface {
	Name() string
	IsDir() bool
	ModTime() time.Time
}

// Parser abstraction for AST operations
type ASTParser interface {
	ParseDir(path string) ([]Symbol, error)
}

// Hasher abstraction for integrity checks
type ContentHasher interface {
	HashDir(path string) (string, error)
}
