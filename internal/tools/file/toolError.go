package file

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/pardnchiu/agenvoy/internal/utils"
)

type ToolError struct {
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	ToolName  string `json:"tool_name"`
	Args      string `json:"args"`
	Error     string `json:"error"`
}

func SaveToolError(sessionID, toolName, args, errMsg string) string {
	raw := toolName + "|" + args + "|" + errMsg
	h := sha256.Sum256([]byte(raw))
	hash := hex.EncodeToString(h[:])[:8]

	configDir, err := utils.GetConfigDir("sessions")
	if err != nil {
		return hash
	}

	record := ToolError{
		Hash:      hash,
		Timestamp: time.Now().Unix(),
		ToolName:  toolName,
		Args:      args,
		Error:     errMsg,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return hash
	}

	now := time.Now()
	dir := filepath.Join(configDir.Home, sessionID, "tool_errors", now.Format("2006-01-02"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return hash
	}
	utils.WriteFile(filepath.Join(dir, hash+".json"), string(data), 0644)
	return hash
}

func GetToolError(sessionID, hash string) string {
	configDir, err := utils.GetConfigDir("sessions")
	if err != nil {
		return ""
	}

	// * search across all date dirs
	sessionDir := filepath.Join(configDir.Home, sessionID, "tool_errors")
	entries, err := os.ReadDir(sessionDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(sessionDir, entry.Name(), hash+".json")
		if data, err := os.ReadFile(path); err == nil {
			return string(data)
		}
	}
	return ""
}
