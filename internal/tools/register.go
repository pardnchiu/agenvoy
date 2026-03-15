package tools

import (
	"context"
	"encoding/json"
	"fmt"

	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"

	_ "github.com/pardnchiu/agenvoy/internal/tools/apis"
	_ "github.com/pardnchiu/agenvoy/internal/tools/apis/googleRSS"
	_ "github.com/pardnchiu/agenvoy/internal/tools/apis/searchWeb"
	_ "github.com/pardnchiu/agenvoy/internal/tools/browser"
	_ "github.com/pardnchiu/agenvoy/internal/tools/calculator"
	_ "github.com/pardnchiu/agenvoy/internal/tools/file"
	_ "github.com/pardnchiu/agenvoy/internal/tools/schedulerTools"
)

func init() {
	toolRegister.Register("run_command", func(ctx context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Command string `json:"command"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return runCommand(ctx, e, params.Command)
	})

	toolRegister.Register("list_tools", func(_ context.Context, e *toolTypes.Executor, _ json.RawMessage) (string, error) {
		type entry struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}

		list := make([]entry, 0, len(e.Tools))
		for _, t := range e.Tools {
			list = append(list, entry{
				Name:        t.Function.Name,
				Description: t.Function.Description,
			})
		}

		if e.APIToolbox != nil {
			for _, raw := range e.APIToolbox.GetTools() {
				fn, _ := raw["function"].(map[string]any)
				if fn == nil {
					continue
				}
				name, _ := fn["name"].(string)
				desc, _ := fn["description"].(string)
				if name != "" {
					list = append(list, entry{Name: name, Description: desc})
				}
			}
		}

		out, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			return "", fmt.Errorf("json.MarshalIndent: %w", err)
		}
		return string(out), nil
	})
}
