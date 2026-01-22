package system

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Josepavese/asdp/engine/domain"
)

type RealFileSystem struct{}

func NewRealFileSystem() *RealFileSystem {
	return &RealFileSystem{}
}

func (fs *RealFileSystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (fs *RealFileSystem) ReadDir(path string) ([]domain.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var infos []domain.FileInfo
	for _, f := range files {
		infos = append(infos, &realFileInfo{info: f})
	}
	return infos, nil
}

func (fs *RealFileSystem) WriteFile(path string, data []byte) error {
	// Ensure parent dir exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

func (fs *RealFileSystem) MkdirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

func (fs *RealFileSystem) Walk(root string, fn func(path string, isDir bool) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return fn(path, info.IsDir())
	})
}

func (fs *RealFileSystem) Stat(path string) (domain.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &realFileInfo{info: info}, nil
}

// Wrapper to match domain interface
type realFileInfo struct {
	info os.FileInfo
}

func (fi *realFileInfo) Name() string {
	return fi.info.Name()
}

func (fi *realFileInfo) IsDir() bool {
	return fi.info.IsDir()
}

func (fi *realFileInfo) ModTime() time.Time {
	return fi.info.ModTime()
}
