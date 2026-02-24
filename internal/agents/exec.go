package agents

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	atypes "github.com/pardnchiu/go-agent-skills/internal/agents/types"
	"github.com/pardnchiu/go-agent-skills/internal/skill"
	"github.com/pardnchiu/go-agent-skills/internal/tools"
	ttypes "github.com/pardnchiu/go-agent-skills/internal/tools/types"
	"github.com/pardnchiu/go-agent-skills/internal/utils"
)

//go:embed prompt/sysPrompt.md
var sysPrompt string

//go:embed prompt/summaryPrompt.md
var summaryPrompt string

//go:embed prompt/sysPromptBase.md
var sysPromptBase string

//go:embed prompt/skillSelector.md
var promptSkillSelectpr string

var (
	MaxToolIterations = 32
)

type Message struct {
	Role       string           `json:"role"`
	Content    any              `json:"content,omitempty"`
	ToolCalls  []OpenAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
}

type OpenAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type OpenAIOutput struct {
	Choices []struct {
		Message      Message `json:"message"`
		Delta        Message `json:"delta"`
		FinishReason string  `json:"finish_reason,omitempty"`
	} `json:"choices"`
	Error *struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Code    json.Number `json:"code"`
	} `json:"error,omitempty"`
}

type Agent interface {
	Send(ctx context.Context, messages []Message, toolDefs []ttypes.Tool) (*OpenAIOutput, error)
	Execute(ctx context.Context, skill *skill.Skill, userInput string, events chan<- atypes.Event, allowAll bool) error
}

func ExecuteAuto(ctx context.Context, agent Agent, scanner *skill.Scanner, userInput string, events chan<- atypes.Event, allowAll bool) error {
	workDir, _ := os.Getwd()

	matched := selectSkill(ctx, agent, scanner, userInput)
	if matched != nil {
		events <- atypes.Event{
			Type: atypes.EventText,
			Text: fmt.Sprintf("Auto-selected skill: %s", matched.Name),
		}
		return Execute(ctx, agent, workDir, matched, userInput, events, allowAll)
	}

	return Execute(ctx, agent, workDir, nil, userInput, events, allowAll)
}

func selectSkill(ctx context.Context, agent Agent, scanner *skill.Scanner, userInput string) *skill.Skill {
	skills := scanner.List()
	if len(skills) == 0 {
		return nil
	}

	var sb strings.Builder
	for _, skill := range skills {
		s := scanner.Skills.ByName[skill]
		sb.WriteString(fmt.Sprintf("- %s: %s\n", skill, s.Description))
	}

	messages := []Message{
		{
			Role:    "system",
			Content: promptSkillSelectpr,
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Available skills:\n%s\nUser request: %s", sb.String(), userInput),
		},
	}

	resp, err := agent.Send(ctx, messages, nil)
	if err != nil || len(resp.Choices) == 0 {
		return nil
	}

	answer := ""
	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		answer = strings.TrimSpace(content)
	}

	if answer == "NONE" || answer == "" {
		return nil
	}

	if s, ok := scanner.Skills.ByName[answer]; ok {
		return s
	}

	cleaned := strings.Trim(answer, "\"'` \n")
	if s, ok := scanner.Skills.ByName[cleaned]; ok {
		return s
	}

	return nil
}

func Execute(ctx context.Context, agent Agent, workDir string, skill *skill.Skill, userInput string, events chan<- atypes.Event, allowAll bool) error {
	if skill != nil && skill.Content == "" {
		return fmt.Errorf("SKILL.md is empty: %s", skill.Path)
	}

	exec, err := tools.NewExecutor(workDir)
	if err != nil {
		return fmt.Errorf("tools.NewExecutor: %w", err)
	}
	now := time.Now()
	dateStr := now.Format("2006-01-02T15:04:05 MST (UTC-07:00)")

	systemPrompt := systemPrompt(workDir, skill)
	toolActions := []Message{}
	history := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("當前時間：%s\n%s", dateStr, userInput)},
	}
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("當前時間：%s\n%s", dateStr, userInput)},
	}

	configDir, err := utils.ConfigDir("sessions")

	indexPath := filepath.Join(configDir.Home, "index.json")
	var sessionID string
	if data, err := os.ReadFile(indexPath); err == nil {
		var indexData struct {
			SessionID string `json:"session_id"`
		}
		if err := json.Unmarshal(data, &indexData); err == nil {
			sessionID = indexData.SessionID

			var summary string
			data, err := os.ReadFile(filepath.Join(configDir.Home, sessionID, "summary.json"))
			if err == nil {
				summary = strings.NewReplacer(
					"{{.Summary}}", string(data),
				).Replace(summaryPrompt)
			}

			data, err = os.ReadFile(filepath.Join(configDir.Home, sessionID, "history.json"))
			if err == nil {
				var oldHistory []Message
				if err := json.Unmarshal(data, &oldHistory); err == nil {
					history = oldHistory
				}
				history = append(history, Message{Role: "user", Content: fmt.Sprintf("當前時間：%s\n%s", dateStr, userInput)})

				recentHistory := oldHistory
				if len(recentHistory) > 4 {
					recentHistory = recentHistory[len(recentHistory)-4:]
				}

				messages = []Message{
					{Role: "system", Content: systemPrompt},
				}
				messages = append(messages, Message{Role: "system", Content: summary})
				messages = append(messages, recentHistory...)
				messages = append(messages, Message{Role: "user", Content: fmt.Sprintf("當前時間：%s\n%s", dateStr, userInput)})
			}
		}
	} else {
		sessionID = fmt.Sprintf("%d", time.Now().UnixNano())
		indexData := struct {
			SessionID string `json:"session_id"`
		}{SessionID: sessionID}

		indexBytes, err := json.Marshal(indexData)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}
		if err := os.WriteFile(indexPath, indexBytes, 0644); err != nil {
			return fmt.Errorf("os.WriteFile: %w", err)
		}
	}

	err = os.MkdirAll(filepath.Join(configDir.Home, sessionID), 0755)
	if err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	for i := 0; i < MaxToolIterations; i++ {
		if i > 0 {
			time.Sleep(1 * time.Second)
		}

		resp, err := agent.Send(ctx, messages, exec.Tools)
		if err != nil {
			return err
		}

		if len(resp.Choices) == 0 {
			return fmt.Errorf("no choices in response")
		}

		choice := resp.Choices[0]

		if len(choice.Message.ToolCalls) > 0 {
			messages = append(messages, choice.Message)

			for _, e := range choice.Message.ToolCalls {
				toolName := e.Function.Name
				if idx := strings.Index(toolName, "<|"); idx != -1 {
					toolName = toolName[:idx]
				}

				events <- atypes.Event{
					Type:     atypes.EventToolCall,
					ToolName: toolName,
					ToolArgs: e.Function.Arguments,
					ToolID:   e.ID,
				}

				if !allowAll {
					replyCh := make(chan bool, 1)
					events <- atypes.Event{
						Type:     atypes.EventToolConfirm,
						ToolName: toolName,
						ToolArgs: e.Function.Arguments,
						ToolID:   e.ID,
						ReplyCh:  replyCh,
					}
					proceed := <-replyCh
					if !proceed {
						events <- atypes.Event{
							Type:     atypes.EventToolSkipped,
							ToolName: toolName,
							ToolID:   e.ID,
						}
						toolActions = append(toolActions, Message{
							Role:       "tool",
							Content:    fmt.Sprintf("Tool '%s' execution skipped by user.", toolName),
							ToolCallID: e.ID,
						})
						messages = append(messages, Message{
							Role:       "tool",
							Content:    "Tool execution skipped by user.",
							ToolCallID: e.ID,
						})
						continue
					}
				}

				result, err := tools.Execute(ctx, exec, toolName, json.RawMessage(e.Function.Arguments))
				if err != nil {
					result = "Error: " + err.Error()
				}

				events <- atypes.Event{
					Type:     atypes.EventToolResult,
					ToolName: toolName,
					ToolID:   e.ID,
					Result:   result,
				}
				toolActions = append(toolActions, Message{
					Role:       "tool",
					Content:    fmt.Sprintf("Tool '%s' executed with result: %s", toolName, result),
					ToolCallID: e.ID,
				})
				messages = append(messages, Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: e.ID,
				})
			}
			continue
		}

		switch v := choice.Message.Content.(type) {
		case string:
			if v != "" {
				const summaryStart = "<!--SUMMARY_START-->"
				const summaryEnd = "<!--SUMMARY_END-->"

				var jsonData any
				var cleaned string

				start := strings.Index(v, summaryStart)
				end := strings.Index(v, summaryEnd)
				if start != -1 && end != -1 && end > start {
					jsonPart := strings.TrimSpace(v[start+len(summaryStart) : end])
					json.Unmarshal([]byte(jsonPart), &jsonData)
					cleaned = strings.TrimRight(v[:start], " \t\n\r")
				} else {
					cleaned = v
				}

				if jsonData != nil {
					path := filepath.Join(configDir.Home, sessionID, "summary.json")
					data, err := json.Marshal(jsonData)
					if err == nil {
						os.WriteFile(path, data, 0644)
					}
				}

				now := time.Now()
				dateStr := now.Format("2006-01-02T15:04:05 MST (UTC-07:00)")
				choice.Message.Content = fmt.Sprintf("當前時間：%s\n%s", dateStr, cleaned)
				messages = append(messages, choice.Message)
				history = append(history, choice.Message)

				events <- atypes.Event{Type: atypes.EventText, Text: cleaned}
			}
		case nil:
		default:
			return fmt.Errorf("unexpected content type: %T", choice.Message.Content)
		}

		filtered := make([]Message, 0, len(history))
		for i, m := range history {
			if m.Role == "system" && i > 0 {
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

		if len(toolActions) > 0 {
			toolActionsDir := filepath.Join(configDir.Work, sessionID)
			if err := os.MkdirAll(toolActionsDir, 0755); err == nil {
				filename := time.Now().Format("2006-01-02-03-04-05") + ".json"
				toolActionsPath := filepath.Join(toolActionsDir, filename)
				if data, err := json.Marshal(toolActions); err == nil {
					os.WriteFile(toolActionsPath, data, 0644)
				}
			}
		}

		events <- atypes.Event{Type: atypes.EventDone}
		return nil
	}

	return fmt.Errorf("exceeded max iterations (%d)", MaxToolIterations)
}

func systemPrompt(workPath string, skill *skill.Skill) string {
	if skill == nil {
		return strings.NewReplacer(
			"{{.WorkPath}}", workPath,
			"{{.Content}}", "",
		).Replace(sysPromptBase)
	}
	content := skill.Content

	for _, prefix := range []string{"scripts/", "templates/", "assets/"} {
		resolved := filepath.Join(skill.Path, prefix)

		if _, err := os.Stat(resolved); err == nil {
			content = strings.ReplaceAll(content, prefix, resolved+string(filepath.Separator))
		}
	}

	return strings.NewReplacer(
		"{{.WorkPath}}", workPath,
		"{{.SkillPath}}", skill.Path,
		"{{.Content}}", content,
	).Replace(sysPrompt)
}
