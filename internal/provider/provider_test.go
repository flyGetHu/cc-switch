package provider

import (
	"testing"
)

func TestGetPresets(t *testing.T) {
	presets := GetPresets()

	if len(presets) != 2 {
		t.Errorf("Expected 2 presets, got %d", len(presets))
	}

	tests := []struct {
		key            string
		name           string
		baseURL        string
		defaultOpus    string
		defaultSonnet  string
		defaultHaiku   string
		smallFast      string
		defaultModel   string
	}{
		{
			key:            "zhipu",
			name:           "智普 GLM",
			baseURL:        "https://open.bigmodel.cn/api/anthropic",
			defaultOpus:    "glm-5",
			defaultSonnet:  "glm-5",
			defaultHaiku:   "glm-5",
			smallFast:      "glm-5",
			defaultModel:   "glm-5",
		},
		{
			key:            "minimax",
			name:           "MiniMax",
			baseURL:        "https://api.minimaxi.com/anthropic",
			defaultOpus:    "MiniMax-M2.5",
			defaultSonnet:  "MiniMax-M2.5",
			defaultHaiku:   "MiniMax-M2.5",
			smallFast:      "MiniMax-M2.5",
			defaultModel:   "MiniMax-M2.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			p, ok := presets[tt.key]
			if !ok {
				t.Fatalf("Preset %s not found", tt.key)
			}

			if p.Name != tt.name {
				t.Errorf("Expected name %s, got %s", tt.name, p.Name)
			}
			if p.BaseURL != tt.baseURL {
				t.Errorf("Expected baseURL %s, got %s", tt.baseURL, p.BaseURL)
			}
			if p.Models.DefaultOpus != tt.defaultOpus {
				t.Errorf("Expected DefaultOpus %s, got %s", tt.defaultOpus, p.Models.DefaultOpus)
			}
			if p.Models.DefaultSonnet != tt.defaultSonnet {
				t.Errorf("Expected DefaultSonnet %s, got %s", tt.defaultSonnet, p.Models.DefaultSonnet)
			}
			if p.Models.DefaultHaiku != tt.defaultHaiku {
				t.Errorf("Expected DefaultHaiku %s, got %s", tt.defaultHaiku, p.Models.DefaultHaiku)
			}
			if p.Models.SmallFast != tt.smallFast {
				t.Errorf("Expected SmallFast %s, got %s", tt.smallFast, p.Models.SmallFast)
			}
			if p.Models.DefaultModel != tt.defaultModel {
				t.Errorf("Expected DefaultModel %s, got %s", tt.defaultModel, p.Models.DefaultModel)
			}
		})
	}
}

func TestProviderStruct(t *testing.T) {
	p := Provider{
		Name:    "Test Provider",
		BaseURL: "https://api.test.com/anthropic",
		Models: ModelConfig{
			DefaultOpus:   "test-opus",
			DefaultSonnet: "test-sonnet",
			DefaultHaiku:  "test-haiku",
			SmallFast:     "test-small-fast",
			DefaultModel:  "test-default",
		},
		APIKey: "test-key",
	}

	if p.Name != "Test Provider" {
		t.Errorf("Expected Name 'Test Provider', got %s", p.Name)
	}
	if p.BaseURL != "https://api.test.com/anthropic" {
		t.Errorf("Expected BaseURL 'https://api.test.com/anthropic', got %s", p.BaseURL)
	}
	if p.APIKey != "test-key" {
		t.Errorf("Expected APIKey 'test-key', got %s", p.APIKey)
	}
}

func TestModelConfigStruct(t *testing.T) {
	mc := ModelConfig{
		DefaultOpus:   "opus-model",
		DefaultSonnet: "sonnet-model",
		DefaultHaiku:  "haiku-model",
		SmallFast:     "small-fast-model",
		DefaultModel:  "default-model",
	}

	if mc.DefaultOpus != "opus-model" {
		t.Errorf("Expected DefaultOpus 'opus-model', got %s", mc.DefaultOpus)
	}
	if mc.DefaultSonnet != "sonnet-model" {
		t.Errorf("Expected DefaultSonnet 'sonnet-model', got %s", mc.DefaultSonnet)
	}
	if mc.DefaultHaiku != "haiku-model" {
		t.Errorf("Expected DefaultHaiku 'haiku-model', got %s", mc.DefaultHaiku)
	}
	if mc.SmallFast != "small-fast-model" {
		t.Errorf("Expected SmallFast 'small-fast-model', got %s", mc.SmallFast)
	}
	if mc.DefaultModel != "default-model" {
		t.Errorf("Expected DefaultModel 'default-model', got %s", mc.DefaultModel)
	}
}
