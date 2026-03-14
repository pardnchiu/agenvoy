package schedulerTools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pardnchiu/agenvoy/internal/filesystem/sessionManager"
	"github.com/pardnchiu/agenvoy/internal/scheduler"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func TaskRoutes(e *toolTypes.Executor, name string, args json.RawMessage) (string, error) {
	switch name {
	case "add_task":
		var params struct {
			At        string `json:"at"`
			Script    string `json:"script"`
			ChannelID string `json:"channel_id"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		if params.ChannelID == "" {
			channelID, err := sessionManager.GetChannelID(e.SessionID)
			if err != nil {
				return "", fmt.Errorf("GetChannelID: %w", err)
			}
			params.ChannelID = channelID
		}
		mgr := scheduler.Get()
		if mgr == nil {
			return "", fmt.Errorf("scheduler not initialized")
		}
		return mgr.AddTask(params.At, params.Script, params.ChannelID)

	case "list_tasks":
		mgr := scheduler.Get()
		if mgr == nil {
			return "", fmt.Errorf("scheduler not initialized")
		}
		tasks := mgr.ListTasks()
		if len(tasks) == 0 {
			return "no onetime tasks", nil
		}
		return strings.Join(tasks, "\n"), nil

	case "remove_task":
		var params struct {
			Index int `json:"index"`
		}
		if err := json.Unmarshal(args, &params); err != nil {
			return "", fmt.Errorf("json.Unmarshal: %w", err)
		}
		mgr := scheduler.Get()
		if mgr == nil {
			return "", fmt.Errorf("scheduler not initialized")
		}
		if err := mgr.RemoveTask(params.Index); err != nil {
			return "", err
		}
		return fmt.Sprintf("onetime task #%d removed", params.Index), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
