package provider

import (
	_ "embed"
	"encoding/json"

	"github.com/pardnchiu/agenvoy/configs"
)

type ProviderItem struct {
	Default string               `json:"default"`
	Models  map[string]ModelItem `json:"models"`
}

type ModelItem struct {
	Input         int    `json:"input"`
	Output        int    `json:"output"`
	Description   string `json:"description"`
	NoTemperature bool   `json:"no_temperature,omitempty"`
}

func parse(data []byte) ProviderItem {
	var cfg ProviderItem
	json.Unmarshal(data, &cfg)
	return cfg
}

func providers() map[string]ProviderItem {
	return map[string]ProviderItem{
		"claude":  parse(configs.ClaudeModels),
		"copilot": parse(configs.CopilotModels),
		"gemini":  parse(configs.GeminiModels),
		"nvidia":  parse(configs.NvidiaModels),
		"openai":  parse(configs.OpenaiModels),
	}
}

func Default(provider string) string {
	return providers()[provider].Default
}

func Get(provider, model string) ModelItem {
	cfg, exist := providers()[provider]
	if !exist {
		return ModelItem{Input: 128000, Output: 16384}
	}

	if info, exist := cfg.Models[model]; exist {
		return info
	}

	if info, exist := cfg.Models[cfg.Default]; exist {
		return info
	}
	return ModelItem{Input: 128000, Output: 16384}
}

func Models(provider string) map[string]ModelItem {
	cfg, exist := providers()[provider]
	if !exist {
		return nil
	}
	return cfg.Models
}

func SupportTemperature(providerName, model string) bool {
	return !Get(providerName, model).NoTemperature
}

func InputBytes(provider, model string) int {
	return Get(provider, model).Input * 4
}

func OutputTokens(provider, model string) int {
	return Get(provider, model).Output
}
