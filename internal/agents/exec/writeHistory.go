package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	atypes "github.com/pardnchiu/go-agent-skills/internal/agents/types"
	"github.com/pardnchiu/go-agent-skills/internal/utils"
)

func writeHistory(choice atypes.OutputChoices, configDir *utils.ConfigDirData, input *atypes.AgentSession, sessionID string) error {
	input.Histories = append(input.Histories, choice.Message)

	filtered := make([]atypes.Message, 0, len(input.Histories))
	for _, m := range input.Histories {
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
