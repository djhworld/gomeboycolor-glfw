package saves

import (
	"io"
	"os"
	"path/filepath"
	"fmt"
)

type FileSystemStore struct {
	baseDir string
}

func NewFileSystemStore(baseDir string) *FileSystemStore {
	f := new(FileSystemStore)
	f.baseDir = baseDir
	return f
}

func (f *FileSystemStore) Open(game string) (io.ReadCloser, error) {
	location := filepath.Join(f.baseDir, game+".sav")

	fmt.Println(location)
	return os.Open(location)
}

func (f *FileSystemStore) Create(game string) (io.WriteCloser, error) {
	location := filepath.Join(f.baseDir, game+".sav")
	fmt.Println(location)
	return os.Create(location)
}
