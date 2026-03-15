package file

import (
	"context"
	"encoding/json"
	"fmt"

	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	toolRegister.Register("read_file", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return read(e, params.Path)
	})

	toolRegister.Register("list_files", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("glob_files", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Pattern string `json:"pattern"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return glob(e, params.Pattern)
	})

	toolRegister.Register("search_content", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Pattern     string `json:"pattern"`
			FilePattern string `json:"file_pattern"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return search(e, params.Pattern, params.FilePattern)
	})

	toolRegister.Register("search_history", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Keyword   string `json:"keyword"`
			TimeRange string `json:"time_range"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return searchHistory(e.SessionID, params.Keyword, params.TimeRange)
	})

	toolRegister.Register("write_file", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("write_script", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Name    string `json:"name"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return writeScript(params.Name, params.Content)
	})

	toolRegister.Register("patch_edit", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("get_tool_error", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("remember_error", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("search_errors", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Keyword string `json:"keyword"`
			Limit   int    `json:"limit"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return SearchErrors(params.Keyword, params.Limit)
	})
}
