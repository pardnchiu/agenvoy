package sessionManager

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"syscall"
	"time"

	agentTypes "github.com/pardnchiu/agenvoy/internal/agents/types"
	"github.com/pardnchiu/agenvoy/internal/filesystem"
)

func SaveToToolCall(sessionID, content string) {
	now := time.Now()
	date := now.Format("2006-01-02")
	toolCallsDir := filepath.Join(filesystem.SessionsDir, sessionID, "tool_calls", date)
	if err := os.MkdirAll(toolCallsDir, 0755); err == nil {
		filename := fmt.Sprintf("%s.json", now.Format("2006-01-02-15-04-05"))
		toolActionsPath := filepath.Join(toolCallsDir, filename)
		if err := filesystem.WriteFile(toolActionsPath, content, 0644); err != nil {
			slog.Warn("WriteFile",
				slog.String("error", err.Error()))
		}
	}
}

func CreateSession() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)

	sessionID := h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:]
	err := os.MkdirAll(filepath.Join(filesystem.SessionsDir, sessionID), 0755)
	if err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}
	return sessionID, nil
}

func LockConfig() (func(), error) {
	lockPath := filepath.Join(filesystem.AgenvoyDir, "config.json.lock")
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile: %w", err)
	}

	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		file.Close()
		return nil, fmt.Errorf("syscall.Flock: %w", err)
	}

	return func() {
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		file.Close()
	}, nil
}

func GetDiscordSession(guildID, channelID, userID string) (string, error) {
	if guildID == "" {
		guildID = "dm"
	}
	if channelID == "" {
		channelID = "ch"
	}
	key := fmt.Sprintf("%s_%s_%s", guildID, channelID, userID)
	sum := sha256.Sum256([]byte(key))

	sessionID := hex.EncodeToString(sum[:])
	sessionDir := filepath.Join(filesystem.SessionsDir, sessionID)
	configPath := filepath.Join(sessionDir, "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(sessionDir, 0755); err != nil {
			return "", fmt.Errorf("os.MkdirAll: %w", err)
		}

		configData, err := json.Marshal(map[string]string{
			"guild_id":   guildID,
			"channel_id": channelID,
			"user_id":    userID,
		})
		if err != nil {
			return "", fmt.Errorf("json.Marshal: %w", err)
		}
		if err := filesystem.WriteFile(configPath, string(configData), 0644); err != nil {
			return "", fmt.Errorf("WriteFile: %w", err)
		}
	}

	return sessionID, nil
}

var MaxHistoryMessages = 16

func GetHistory(sessionID string) []agentTypes.Message {
	sessionDir := filepath.Join(filesystem.SessionsDir, sessionID)
	historyPath := filepath.Join(sessionDir, "history.json")

	data, err := os.ReadFile(historyPath)
	if err != nil {
		return nil
	}
	var oldHistory []agentTypes.Message
	if err := json.Unmarshal(data, &oldHistory); err != nil {
		return nil
	}
	if len(oldHistory) > MaxHistoryMessages {
		oldHistory = oldHistory[len(oldHistory)-MaxHistoryMessages:]
	}
	return oldHistory
}

func SaveHistory(sessionID, content string) error {
	sessionDir := filepath.Join(filesystem.SessionsDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	historyPath := filepath.Join(sessionDir, "history.json")
	if err := filesystem.WriteFile(historyPath, content, 0644); err != nil {
		return fmt.Errorf("WriteFile: %w", err)
	}
	return nil
}
