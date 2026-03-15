package googleRSS

import (
	"context"
	"encoding/json"
	"fmt"

	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	toolRegister.Register("fetch_google_rss", func(ctx context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Keyword string `json:"keyword"`
			Time    string `json:"time"`
			Lang    string `json:"lang"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return Fetch(params.Keyword, params.Time, params.Lang)
	})
}
