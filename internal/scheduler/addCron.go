package scheduler

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pardnchiu/agenvoy/internal/filesystem"
)

type cronItem struct {
	line       string
	expression string
	script     string
	channelID  string
	cronID     int64
}

func (s *Scheduler) AddCron(expression, script, channelID string) error {
	line := buildLine(expression, script, channelID)
	item, err := parseCronLine(line)
	if err != nil {
		return err
	}

	id, err := s.cron.Add(item.expression, s.makeCronAction(item))
	if err != nil {
		return fmt.Errorf("cron.Add: %w", err)
	}
	item.cronID = id

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := appendLine(filesystem.CronsPath, line); err != nil {
		s.cron.Remove(id)
		return fmt.Errorf("appendLine: %w", err)
	}

	s.crons = append(s.crons, item)
	return nil
}

func parseCronLine(line string) (cronItem, error) {
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return cronItem{}, fmt.Errorf("at least 6 fields `{min} {hour} {dom} {mon} {dow} {script}`")
	}

	channelID := ""
	if len(fields) >= 7 {
		channelID = fields[6]
	}

	return cronItem{
		line:       line,
		expression: strings.Join(fields[:5], " "),
		script:     fields[5],
		channelID:  channelID,
	}, nil
}

func (s *Scheduler) makeCronAction(item cronItem) func() {
	return func() {
		output := runScript("cron", filepath.Join(filesystem.ScriptsDir, item.script))
		s.mu.Lock()
		cb := s.OnCompleted
		s.mu.Unlock()
		if item.channelID != "" && cb != nil {
			cb(item.channelID, output)
		}
	}
}
