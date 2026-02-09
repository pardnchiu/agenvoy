package copilot

import (
	"context"
	"fmt"
	"io"

	"github.com/pardnchiu/go-agent-skills/internal/agents"
	"github.com/pardnchiu/go-agent-skills/internal/skill"
	"github.com/pardnchiu/go-agent-skills/internal/tools"
	"github.com/pardnchiu/go-agent-skills/internal/utils"
)

var (
	DefaultModel = "gpt-4.1"
	ChatAPI      = "https://api.githubcopilot.com/chat/completions"
)

func (a *Agent) Execute(ctx context.Context, skill *skill.Skill, userInput string, output io.Writer, allowAll bool) error {
	if err := a.checkExpires(ctx); err != nil {
		return err
	}
	return agents.Execute(ctx, a, skill, userInput, output, allowAll)
}

func (a *Agent) SendChat(ctx context.Context, messages []agents.Message, toolDefs []tools.Tool) (*agents.OpenAIOutput, error) {
	result, _, err := utils.POSTJson[agents.OpenAIOutput](ctx, a.httpClient, ChatAPI, map[string]string{
		"Authorization":         "Bearer " + a.Refresh.Token,
		"Editor-Version":        "vscode/1.95.0",
		"Editor-Plugin-Version": "copilot/1.245.0",
		"Openai-Organization":   "github-copilot",
	}, map[string]any{
		"model":    DefaultModel,
		"messages": messages,
		"tools":    toolDefs,
	})
	if err != nil {
		return nil, fmt.Errorf("API request: %w", err)
	}
	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Message)
	}

	return &result, nil
}

func (a *Agent) GetWorkDir() string {
	return a.workDir
}
