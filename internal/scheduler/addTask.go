package scheduler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pardnchiu/agenvoy/internal/filesystem"
)

type taskItem struct {
	line      string
	at        time.Time
	script    string
	channelID string
}

// * allow: +5m, +1h30m, 15:04, 2006-01-02 15:04, RFC3339
func (s *Scheduler) AddTask(text, script, channelID string) (string, error) {
	at, err := parseTaskTime(text)
	if err != nil {
		return "", err
	}

	if !at.After(time.Now()) {
		return "", fmt.Errorf("already gone")
	}

	line := buildLine(at.UTC().Format(time.RFC3339), script, channelID)
	item, err := parseTaskLine(line)
	// * ensure format is correct
	if err != nil {
		return "", fmt.Errorf("parseTaskLine: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := appendLine(filesystem.TasksPath, line); err != nil {
		return "", fmt.Errorf("appendLine: %w", err)
	}

	if err := s.setTask(item); err != nil {
		return "", fmt.Errorf("s.setTask: %w", err)
	}
	return fmt.Sprintf("already set up task at %s: %s", at.Local().Format("2006-01-02 15:04:05"), script), nil
}

func parseTaskTime(text string) (time.Time, error) {
	text = strings.TrimSpace(text)

	if strings.HasPrefix(text, "+") {
		duration, err := time.ParseDuration(text[1:])
		if err != nil {
			return time.Time{}, fmt.Errorf("time.ParseDuration: %w", err)
		}
		return time.Now().Add(duration), nil
	}

	if t, err := time.ParseInLocation("2006-01-02 15:04", text, time.Local); err == nil {
		return t, nil
	}

	if t, err := time.Parse(time.RFC3339, text); err == nil {
		return t, nil
	}

	// * 15:04, no date, assume today
	if t, err := time.ParseInLocation("15:04", text, time.Local); err == nil {
		now := time.Now()
		result := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		if !result.After(now) {
			return time.Time{}, fmt.Errorf("already gone: %q", text)
		}
		return result, nil
	}

	return time.Time{}, fmt.Errorf("parseTime: %s", text)
}

func parseTaskLine(line string) (taskItem, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return taskItem{}, fmt.Errorf("at least 2 fields `{time} {script}`")
	}

	at, err := time.Parse(time.RFC3339, fields[0])
	if err != nil {
		return taskItem{}, fmt.Errorf("not RFC3339: %w", err)
	}

	item := taskItem{
		line:   line,
		at:     at,
		script: fields[1],
	}
	if len(fields) >= 3 {
		item.channelID = fields[2]
	}
	return item, nil
}

func (s *Scheduler) setTask(item taskItem) error {
	if err := os.MkdirAll(filesystem.ScriptsDir, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	scriptPath := filepath.Join(filesystem.ScriptsDir, item.script)
	delay := time.Until(item.at)

	execTime := time.AfterFunc(delay, func() {
		output := runScript("task", scriptPath)

		if item.channelID != "" {
			s.mu.Lock()
			cb := s.OnCompleted
			s.mu.Unlock()

			if cb != nil {
				cb(item.channelID, output)
			}
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		delete(s.timers, item.line)
		removeLine(filesystem.TasksPath, item.line)
		// * remove in memory cache, ensure completed remove
		s.removeTask(item.line)
		removeScript(scriptPath)
	})
	s.timers[item.line] = execTime
	s.tasks = append(s.tasks, item)
	return nil
}

func (s *Scheduler) removeTask(line string) {
	for i, task := range s.tasks {
		if task.line == line {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return
		}
	}
}
