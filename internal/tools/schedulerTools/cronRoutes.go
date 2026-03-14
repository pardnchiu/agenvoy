package schedulerTools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pardnchiu/agenvoy/internal/filesystem/sessionManager"
	"github.com/pardnchiu/agenvoy/internal/scheduler"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func Routes(e *toolTypes.Executor, name string, args json.RawMessage) (string, error) {
	switch name {
	case "add_cron":
		var params struct {
			CronExpr  string `json:"cron_expr"`
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
		if err := mgr.AddCron(params.CronExpr, params.Script, params.ChannelID); err != nil {
			return "", err
		}
		return fmt.Sprintf("cron task added: %s %s", params.CronExpr, params.Script), nil

	case "list_crons":
		mgr := scheduler.Get()
		if mgr == nil {
			return "", fmt.Errorf("scheduler not initialized")
		}
		tasks := mgr.ListCronTasks()
		if len(tasks) == 0 {
			return "no cron tasks", nil
		}
		return strings.Join(tasks, "\n"), nil

	case "remove_cron":
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
		if err := mgr.RemoveCronTask(params.Index); err != nil {
			return "", err
		}
		return fmt.Sprintf("cron task #%d removed", params.Index), nil

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}
