package schedulerTools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pardnchiu/agenvoy/internal/filesystem/sessionManager"
	"github.com/pardnchiu/agenvoy/internal/scheduler"
	toolRegister "github.com/pardnchiu/agenvoy/internal/tools/register"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func init() {
	toolRegister.Register("add_cron", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("list_crons", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		mgr := scheduler.Get()
		if mgr == nil {
			return "", fmt.Errorf("scheduler not initialized")
		}
		tasks := mgr.ListCronTasks()
		if len(tasks) == 0 {
			return "no cron tasks", nil
		}
		return strings.Join(tasks, "\n"), nil
	})

	toolRegister.Register("remove_cron", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("add_task", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})

	toolRegister.Register("list_tasks", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
		mgr := scheduler.Get()
		if mgr == nil {
			return "", fmt.Errorf("scheduler not initialized")
		}
		tasks := mgr.ListTasks()
		if len(tasks) == 0 {
			return "no onetime tasks", nil
		}
		return strings.Join(tasks, "\n"), nil
	})

	toolRegister.Register("remove_task", func(_ context.Context, e *toolTypes.Executor, args json.RawMessage) (string, error) {
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
	})
}
