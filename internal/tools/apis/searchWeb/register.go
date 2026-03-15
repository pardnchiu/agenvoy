package searchWeb

import (
	"context"
	"encoding/json"
	"fmt"

	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	toolRegister.Register("search_web", func(ctx context.Context, _ *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Query string `json:"query"`
			Range string `json:"range"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return Search(ctx, params.Query, TimeRange(params.Range))
	})
}
