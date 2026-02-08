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
	WorkPath       string
	Allowed        []string // limit to these folders to use
	AllowedCommand map[string]bool
	Exclude        []string
	Tools          []Tool
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
		AllowedCommand: map[string]bool{
			// Version Control
			"git": true,

			// Programming Languages & Package Managers
			"go":      true,
			"node":    true,
			"npm":     true,
			"yarn":    true,
			"pnpm":    true,
			"python":  true,
			"python3": true,
			"pip":     true,
			"pip3":    true,

			// File Operations
			"ls":    true,
			"cat":   true,
			"head":  true,
			"tail":  true,
			"pwd":   true,
			"mkdir": true,
			"touch": true,
			"cp":    true,
			"mv":    true,
			"rm":    true, // * not support native rm, but move to .Trash

			// Text Processing
			"grep": true,
			"sed":  true,
			"awk":  true,
			"sort": true,
			"uniq": true,
			"diff": true,
			"cut":  true,
			"tr":   true,
			"wc":   true,

			// Search & Find
			"find": true,

			// Data Format
			"jq": true,

			// System Info
			"echo":  true,
			"which": true,
			"date":  true,
		},
		Exclude: []string{
			".DS_Store", ".git", "node_modules", "vendor", ".vscode", ".idea", "dist", "build",
		},
		Tools: tools,
	}, nil
}

func (e *Executor) Execute(name string, args json.RawMessage) (string, error) {
	switch name {
	case "read_file":
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("failed to unmarshal json (%s): %w", name, err)
		}
		return e.readFile(params.Path)

	case "list_files":
		var params struct {
			Path      string `json:"path"`
			Recursive bool   `json:"recursive"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("failed to unmarshal json (%s): %w", name, err)
		}
		return e.listFiles(params.Path, params.Recursive)

	case "glob_files":
		var params struct {
			Pattern string `json:"pattern"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("failed to unmarshal json (%s): %w", name, err)
		}
		return e.globFiles(params.Pattern)

	case "write_file":
		var params struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("failed to unmarshal json (%s): %w", name, err)
		}
		return e.writeFile(params.Path, params.Content)

	case "search_content":
		var params struct {
			Pattern     string `json:"pattern"`
			FilePattern string `json:"file_pattern"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", err
		}
		return e.searchContent(params.Pattern, params.FilePattern)

	case "run_command":
		var params struct {
			Command string `json:"command"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", err
		}
		return e.runCommand(params.Command)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
