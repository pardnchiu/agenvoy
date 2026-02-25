package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pardnchiu/go-agent-skills/internal/utils"
)

func writeHistory(choice OpenAIOutputChoices, configDir *utils.ConfigDirData, input *SessionData, sessionID string) error {
	input.histories = append(input.histories, choice.Message)

	filtered := make([]Message, 0, len(input.histories))
	for _, m := range input.histories {
		if m.Role == "system" {
			continue
		}
		if m.Role == "assistant" && len(m.ToolCalls) > 0 {
			continue
		}
		if m.Role == "tool" {
			continue
		}
		filtered = append(filtered, m)
	}

	historyPath := filepath.Join(configDir.Home, sessionID, "history.json")
	historyData, err := json.Marshal(filtered)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	if err := os.WriteFile(historyPath, historyData, 0644); err != nil {
		return fmt.Errorf("os.WriteFile: %w", err)
	}
	return nil
}
