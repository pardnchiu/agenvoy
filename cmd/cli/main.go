package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sort"

	copilot "github.com/pardnchiu/go-agent-skills/internal/agents/copilot"
	"github.com/pardnchiu/go-agent-skills/internal/skill"
)

func main() {
	client, err := copilot.New()
	if err != nil {
		slog.Error("failed to load Copilot token",
			slog.String("error", err.Error()))
		os.Exit(1)
	}

	if os.Args[1] == "list" {
		scanner := skill.NewScanner()

		if len(scanner.Skills.ByName) == 0 {
			fmt.Println("No skills found")
			fmt.Println("\nScanned paths:")
			for _, path := range scanner.Skills.Paths {
				fmt.Printf("  - %s\n", path)
			}
			return
		}

		names := scanner.List()
		sort.Strings(names)

		fmt.Printf("Found %d skill(s):\n\n", len(names))
		for _, name := range names {
			s := scanner.Skills.ByName[name]
			fmt.Printf("• %s\n", name)
			if s.Description != "" {
				fmt.Printf("  %s\n", s.Description)
			}
			fmt.Printf("  Path: %s\n\n", s.Path)
		}
		return
	}

	if os.Args[1] == "run" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: go run cmd/cli/main.go run <skill_name> <input>")
			os.Exit(1)
		}

		skillName := os.Args[2]
		userInput := os.Args[3]
		allowAll := false

		// Check for --allow flag
		if slices.Contains(os.Args[4:], "--allow") {
			allowAll = true
		}

		client, err := copilot.New()
		if err != nil {
			slog.Error("failed to load Copilot token", slog.String("error", err.Error()))
			os.Exit(1)
		}

		scanner := skill.NewScanner()
		targetSkill, ok := scanner.Skills.ByName[skillName]
		if !ok {
			slog.Error("skill not found", slog.String("name", skillName))
			os.Exit(1)
		}

		ctx := context.Background()
		if err := client.Execute(ctx, targetSkill, userInput, os.Stdout, allowAll); err != nil {
			slog.Error("failed to execute skill", slog.String("error", err.Error()))
			os.Exit(1)
		}
		return
	}

	slog.Info("successfully loaded Copilot token",
		slog.String("access_token", client.Token.AccessToken),
		slog.String("token_type", client.Token.TokenType),
		slog.String("scope", client.Token.Scope),
		slog.Time("expires_at", client.Token.ExpiresAt))

	// if len(os.Args) < 3 || os.Args[1] != "input" {
	// 	slog.Error("usage: go run cmd/cli/main.go input \"your message\"")
	// 	os.Exit(1)
	// }

	// userInput := os.Args[2]
	// ctx := context.Background()
	// messages := []c.Message{
	// 	{
	// 		Role:    "system",
	// 		Content: "You are a helpful assistant.",
	// 	},
	// 	{
	// 		Role:    "user",
	// 		Content: userInput,
	// 	},
	// }

	// resp, err := client.SendChat(ctx, messages, nil)
	// if err != nil {
	// 	slog.Error("failed to send chat",
	// 		slog.String("error", err.Error()))
	// 	os.Exit(1)
	// }

	// if len(resp.Choices) > 0 {
	// 	choice := resp.Choices[0]
	// 	if content, ok := choice.Message.Content.(string); ok {
	// 		fmt.Println("Response:", content)
	// 	}
	// }
}
