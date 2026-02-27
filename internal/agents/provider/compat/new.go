package compat

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Agent struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	workDir    string
}

var (
	defaultModel = "qwen3:8b"
	prefix       = "compat@"
)

func New(model ...string) (*Agent, error) {
	if len(model) > 0 && strings.HasPrefix(model[0], prefix) {
		defaultModel = strings.TrimPrefix(model[0], prefix)
	}

	baseURL := os.Getenv("COMPAT_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	baseURL = strings.TrimRight(baseURL, "/")

	apiKey := os.Getenv("COMPAT_API_KEY")

	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd: %w", err)
	}

	return &Agent{
		httpClient: &http.Client{},
		baseURL:    baseURL,
		apiKey:     apiKey,
		workDir:    workDir,
	}, nil
}
