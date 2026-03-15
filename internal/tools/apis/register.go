package apis

import (
	"context"
	"encoding/json"
	"fmt"

	apiAdapter "github.com/pardnchiu/agenvoy/internal/tools/apis/adapter"
	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	// * api adapter
	toolRegister.Register("send_http_request", func(ctx context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		var params struct {
			URL         string            `json:"url"`
			Method      string            `json:"method"`
			Headers     map[string]string `json:"headers"`
			Body        map[string]any    `json:"body"`
			ContentType string            `json:"content_type"`
			Timeout     int               `json:"timeout"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return apiAdapter.Send(params.URL, params.Method, params.Headers, params.Body, params.ContentType, params.Timeout)
	})
}
