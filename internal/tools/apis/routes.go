package apis

import (
	"encoding/json"
	"fmt"

	"github.com/pardnchiu/agenvoy/internal/tools/apiAdapter"
	"github.com/pardnchiu/agenvoy/internal/tools/apis/googleRSS"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func Routes(e *toolTypes.Executor, name string, args json.RawMessage) (string, error) {
	switch name {
	case "send_http_request":
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

	case "fetch_google_rss":
		var params struct {
			Keyword string `json:"keyword"`
			Time    string `json:"time"`
			Lang    string `json:"lang"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		return googleRSS.Fetch(params.Keyword, params.Time, params.Lang)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
