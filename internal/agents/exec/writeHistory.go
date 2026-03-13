package exec

import (
	"encoding/json"
	"fmt"

	agentTypes "github.com/pardnchiu/agenvoy/internal/agents/types"
	"github.com/pardnchiu/agenvoy/internal/filesystem/sessionManager"
)

func writeHistory(choice agentTypes.OutputChoices, session *agentTypes.AgentSession) error {
	session.Histories = append(session.Histories, choice.Message)

	filtered := make([]agentTypes.Message, 0, len(session.Histories))
	for _, m := range session.Histories {
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
	historyData, err := json.Marshal(filtered)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = sessionManager.SaveHistory(session.ID, string(historyData))
	if err != nil {
		return fmt.Errorf("sessionManager.SaveHistory: %w", err)
	}
	return nil
}
