package calculator

import (
	"context"
	"encoding/json"
	"fmt"

	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	toolRegister.Register("calculate", func(_ context.Context, _ *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			Expression string `json:"expression"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return Calc(params.Expression)
	})
}
