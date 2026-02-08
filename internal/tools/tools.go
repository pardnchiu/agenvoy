package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func (e *Executor) readFile(path string) (string, error) {
	fullPath := getFullPath(e, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("tools/readFile/{os.ReadFile}: %w", err)
	}
	return string(data), nil
}

func getFullPath(e *Executor, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(e.WorkPath, path)
}
