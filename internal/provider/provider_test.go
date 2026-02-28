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
		key     string
		name    string
		baseURL string
		opus    string
		sonnet  string
		haiku   string
	}{
		{
			key:     "zhipu",
			name:    "智普 GLM",
			baseURL: "https://open.bigmodel.cn/api/anthropic",
			opus:    "glm-5",
			sonnet:  "glm-5",
			haiku:   "glm-5",
		},
		{
			key:     "minimax",
			name:    "MiniMax",
			baseURL: "https://api.minimaxi.com/anthropic",
			opus:    "MiniMax-M2.5",
			sonnet:  "MiniMax-M2.5",
			haiku:   "MiniMax-M2.5",
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
			if p.Models.Opus != tt.opus {
				t.Errorf("Expected Opus %s, got %s", tt.opus, p.Models.Opus)
			}
			if p.Models.Sonnet != tt.sonnet {
				t.Errorf("Expected Sonnet %s, got %s", tt.sonnet, p.Models.Sonnet)
			}
			if p.Models.Haiku != tt.haiku {
				t.Errorf("Expected Haiku %s, got %s", tt.haiku, p.Models.Haiku)
			}
		})
	}
}

func TestProviderStruct(t *testing.T) {
	p := Provider{
		Name:    "Test Provider",
		BaseURL: "https://api.test.com/anthropic",
		Models: ModelConfig{
			Opus:   "test-opus",
			Sonnet: "test-sonnet",
			Haiku:  "test-haiku",
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
		Opus:   "opus-model",
		Sonnet: "sonnet-model",
		Haiku:  "haiku-model",
	}

	if mc.Opus != "opus-model" {
		t.Errorf("Expected Opus 'opus-model', got %s", mc.Opus)
	}
	if mc.Sonnet != "sonnet-model" {
		t.Errorf("Expected Sonnet 'sonnet-model', got %s", mc.Sonnet)
	}
	if mc.Haiku != "haiku-model" {
		t.Errorf("Expected Haiku 'haiku-model', got %s", mc.Haiku)
	}
}
