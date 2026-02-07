package tools

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type Executor struct {
	WorkPath string
	Allowed  []string // limit to these folders to use
	Tools    []Tool
}

//go:embed tools.json
var toolsJSON []byte

func NewExecutor(workPath string) (*Executor, error) {
	var tools []Tool
	if err := json.Unmarshal(toolsJSON, &tools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tools: %w", err)
	}

	return &Executor{
		WorkPath: workPath,
		Tools:    tools,
	}, nil
}

func (e *Executor) Execute(name string, args json.RawMessage) (string, error) {
	switch name {
	case "read_file":
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", err
		}
		return e.readFile(params.Path)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func (e *Executor) getPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(e.WorkPath, path)
}

func (e *Executor) readFile(path string) (string, error) {
	fullPath := e.getPath(path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	content := string(data)
	return content, nil
}
