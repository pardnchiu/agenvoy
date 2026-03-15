package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pardnchiu/agenvoy/internal/filesystem"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

func write(e *toolTypes.Executor, path, content string) (string, error) {
	if content == "" {
		return "", fmt.Errorf("refused to write empty content to file (%s)", path)
	}

	fullPath, err := getFullPath(e, path)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory (%s): %w", path, err)
	}

	if err := filesystem.WriteFile(fullPath, content, 0644); err != nil {
		return "", fmt.Errorf("utils.WriteFile: %w", err)
	}

	return fmt.Sprintf("Successfully wrote file: %s", path), nil
}

func writeScript(name, content string) (string, error) {
	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".sh" && ext != ".py" {
		return "", fmt.Errorf("scripts only support .sh or .py")
	}
	if filepath.Base(name) != name {
		return "", fmt.Errorf("must not contain path separator")
	}

	if err := os.MkdirAll(filesystem.ScriptsDir, 0755); err != nil {
		return "", fmt.Errorf("os.MkdirAll: %w", err)
	}

	base := strings.TrimSuffix(name, ext)
	uniqueName := fmt.Sprintf("%s_%d%s", base, time.Now().UTC().Unix(), ext)

	path := filepath.Join(filesystem.ScriptsDir, uniqueName)
	if err := filesystem.WriteFile(path, content, 0755); err != nil {
		return "", fmt.Errorf("filesystem.WriteFile: %w", err)
	}

	return fmt.Sprintf(`script saved. pass "%s" as the script parameter to add_task or add_cron`, uniqueName), nil
}
