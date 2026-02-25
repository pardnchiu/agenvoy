package exec

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/pardnchiu/go-agent-skills/internal/utils"
)

func extractSummary(configDir *utils.ConfigDirData, sessionID, value string) string {
	const summaryStart = "<!--SUMMARY_START-->"
	const summaryEnd = "<!--SUMMARY_END-->"

	var jsonData any
	var cleaned string

	start := strings.Index(value, summaryStart)
	end := strings.Index(value, summaryEnd)
	if start != -1 && end != -1 && end > start {
		jsonPart := strings.TrimSpace(value[start+len(summaryStart) : end])
		json.Unmarshal([]byte(jsonPart), &jsonData)
		cleaned = strings.TrimRight(value[:start], " \t\n\r")
	} else {
		cleaned = value
	}

	if jsonData != nil {
		path := filepath.Join(configDir.Home, sessionID, "summary.json")
		data, err := json.Marshal(jsonData)
		if err == nil {
			os.WriteFile(path, data, 0644)
		}
	}
	return cleaned
}
