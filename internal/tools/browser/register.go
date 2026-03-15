package browser

import (
	"context"
	"encoding/json"
	"fmt"

	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	toolRegister.Register("fetch_page", func(_ context.Context, _ *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return Load(params.URL)
	})

	toolRegister.Register("download_page", func(_ context.Context, _ *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Href   string `json:"href"`
			SaveTo string `json:"save_to"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return Download(params.Href, params.SaveTo)
	})
}
