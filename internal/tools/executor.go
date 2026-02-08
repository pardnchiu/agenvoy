package tools

import (
	_ "embed"
	"encoding/json"
	"fmt"
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
