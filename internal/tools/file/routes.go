package file

import (
	"encoding/json"
	"fmt"

	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func Routes(e *toolTypes.Executor, name string, args json.RawMessage) (string, error) {
	switch name {
	case "read_file":
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return read(e, params.Path)

	case "list_files":
		var params struct {
			Path      string `json:"path"`
			Recursive bool   `json:"recursive"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		if isDenied(params.Path) {
			return "", fmt.Errorf("access denied: %s", params.Path)
		}
		return list(e, params.Path, params.Recursive)

	case "glob_files":
		var params struct {
			Pattern string `json:"pattern"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return glob(e, params.Pattern)

	case "search_content":
		var params struct {
			Pattern     string `json:"pattern"`
			FilePattern string `json:"file_pattern"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return search(e, params.Pattern, params.FilePattern)

	case "search_history":
		var params struct {
			Keyword   string `json:"keyword"`
			TimeRange string `json:"time_range"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return searchHistory(e.SessionID, params.Keyword, params.TimeRange)

	case "write_file":
		var params struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		if isDenied(params.Path) {
			return "", fmt.Errorf("access denied: %s", params.Path)
		}
		return write(e, params.Path, params.Content)

	case "write_script":
		var params struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return writeScript(params.Name, params.Content)

	case "patch_edit":
		var params struct {
			Path      string `json:"path"`
			OldString string `json:"old_string"`
			NewString string `json:"new_string"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		if isDenied(params.Path) {
			return "", fmt.Errorf("access denied: %s", params.Path)
		}
		return patch(e, params.Path, params.OldString, params.NewString)

	case "get_tool_error":
		var params struct {
			Hash string `json:"hash"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		result := GetToolError(e.SessionID, params.Hash)
		if result == "" {
			return "not found", nil
		}
		return result, nil

	case "remember_error":
		var params struct {
			ToolName string   `json:"tool_name"`
			Keywords []string `json:"keywords"`
			Symptom  string   `json:"symptom"`
			Cause    string   `json:"cause"`
			Action   string   `json:"action"`
			Outcome  string   `json:"outcome"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return SaveErrorMemory(e.SessionID, ErrorMemory{
			ToolName: params.ToolName,
			Keywords: params.Keywords,
			Symptom:  params.Symptom,
			Cause:    params.Cause,
			Action:   params.Action,
			Outcome:  params.Outcome,
		})

	case "search_errors":
		var params struct {
			Keyword string `json:"keyword"`
			Limit   int    `json:"limit"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return SearchErrors(params.Keyword, params.Limit)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
