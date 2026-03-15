package file

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/pardnchiu/agenvoy/configs"
	toolTypes "github.com/pardnchiu/agenvoy/internal/tools/types"
)

type deniedConfig struct {
	Dirs       []string `json:"dirs"`
	Files      []string `json:"files"`
	Prefixes   []string `json:"prefixes"`
	Extensions []string `json:"extensions"`
}

var DeniedConfig = func() deniedConfig {
	var cfg deniedConfig
	if err := json.Unmarshal(configs.DeniedMap, &cfg); err != nil {
		slog.Warn("json.Unmarshal",
			slog.String("error", err.Error()))
	}
	return cfg
}()

func isDenied(path string) bool {
	cleaned := filepath.Clean(path)
	base := filepath.Base(cleaned)

	for _, dir := range DeniedConfig.Dirs {
		if strings.Contains(cleaned, fmt.Sprintf("/%s/", dir)) || strings.Contains(cleaned, fmt.Sprintf("/%s", dir)) {
			return true
		}
	}
	for _, f := range DeniedConfig.Files {
		if strings.Contains(cleaned, f) {
			return true
		}
	}
	for _, prefix := range DeniedConfig.Prefixes {
		if strings.HasPrefix(base, prefix) && !strings.Contains(base, ".example") {
			return true
		}
	}
	for _, ext := range DeniedConfig.Extensions {
		if strings.HasSuffix(base, ext) {
			return true
		}
	}
	return false
}

func read(e *toolTypes.Executor, path string) (string, error) {
	fullPath, err := getFullPath(e, path)
	if err != nil {
		return "", err
	}

	if isDenied(fullPath) {
		return "", fmt.Errorf("access denied: %s", path)
	}

	if isExclude(e, fullPath) {
		return "", fmt.Errorf("path is excluded: %s", path)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file (%s): %w", path, err)
	}
	return string(data), nil
}

func getFullPath(e *toolTypes.Executor, path string) (string, error) {
	if !filepath.IsAbs(path) {
		return filepath.Join(e.WorkPath, path), nil
	}
	cleaned := filepath.Clean(path)
	homeDir, err := os.UserHomeDir()
	if err != nil || !strings.HasPrefix(cleaned, filepath.Clean(homeDir)+string(filepath.Separator)) {
		return "", fmt.Errorf("only allow user home: %s", path)
	}
	return cleaned, nil
}

func isExclude(e *toolTypes.Executor, path string) bool {
	excluded := false
	for _, e := range e.Exclude {
		match, err := filepath.Match(e.File, filepath.Base(path))
		if err != nil {
			continue
		}

		if !match {
			match = strings.Contains(path, "/"+e.File+"/") ||
				strings.HasPrefix(path, e.File+"/")
		}
		if match {
			excluded = !e.Negate
		}
	}
	return excluded
}
