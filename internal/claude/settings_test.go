package claude

import (
	"os"
	"path/filepath"
	"testing"

	"cc-switch/internal/provider"
)

func TestGetSettingsPath(t *testing.T) {
	path := GetSettingsPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".claude", "settings.json")

	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestSetSettingsPath(t *testing.T) {
	testPath := "/tmp/test/settings.json"
	SetSettingsPath(testPath)

	if GetSettingsPath() != testPath {
		t.Errorf("Expected %s, got %s", testPath, GetSettingsPath())
	}

	SetSettingsPath("")
}

func TestReadSettings_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	s, err := ReadSettings()
	if err != nil {
		t.Errorf("ReadSettings failed: %v", err)
	}

	if s.Env == nil {
		t.Error("Env should not be nil")
	}
}

func TestReadAndWriteSettings(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	original := &Settings{
		Env: map[string]string{
			"TEST_VAR": "test_value",
		},
		Permissions: Permissions{
			Allow: []string{"read"},
			Deny:  []string{"write"},
		},
	}

	if err := WriteSettings(original); err != nil {
		t.Fatalf("WriteSettings failed: %v", err)
	}

	loaded, err := ReadSettings()
	if err != nil {
		t.Fatalf("ReadSettings failed: %v", err)
	}

	if loaded.Env["TEST_VAR"] != "test_value" {
		t.Errorf("Expected TEST_VAR 'test_value', got %s", loaded.Env["TEST_VAR"])
	}

	if len(loaded.Permissions.Allow) != 1 || loaded.Permissions.Allow[0] != "read" {
		t.Errorf("Permissions.Allow not preserved correctly")
	}
}

func TestApplyProvider(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	original := &Settings{
		Env: map[string]string{
			"EXISTING_VAR": "value",
		},
		Permissions: Permissions{
			Allow: []string{},
			Deny:  []string{},
		},
	}

	if err := WriteSettings(original); err != nil {
		t.Fatalf("WriteSettings failed: %v", err)
	}

	p := &provider.Provider{
		Name:    "Test",
		BaseURL: "https://api.test.com/anthropic",
		Models: provider.ModelConfig{
			Opus:   "test-opus",
			Sonnet: "test-sonnet",
			Haiku:  "test-haiku",
		},
		APIKey: "test-api-key",
	}

	if err := ApplyProvider(p); err != nil {
		t.Fatalf("ApplyProvider failed: %v", err)
	}

	loaded, err := ReadSettings()
	if err != nil {
		t.Fatalf("ReadSettings failed: %v", err)
	}

	if loaded.Env["ANTHROPIC_BASE_URL"] != p.BaseURL {
		t.Errorf("Expected ANTHROPIC_BASE_URL %s, got %s", p.BaseURL, loaded.Env["ANTHROPIC_BASE_URL"])
	}
	if loaded.Env["ANTHROPIC_AUTH_TOKEN"] != p.APIKey {
		t.Errorf("Expected ANTHROPIC_AUTH_TOKEN %s, got %s", p.APIKey, loaded.Env["ANTHROPIC_AUTH_TOKEN"])
	}
	if loaded.Env["ANTHROPIC_DEFAULT_OPUS_MODEL"] != p.Models.Opus {
		t.Errorf("Expected ANTHROPIC_DEFAULT_OPUS_MODEL %s, got %s", p.Models.Opus, loaded.Env["ANTHROPIC_DEFAULT_OPUS_MODEL"])
	}
	if loaded.Env["ANTHROPIC_DEFAULT_SONNET_MODEL"] != p.Models.Sonnet {
		t.Errorf("Expected ANTHROPIC_DEFAULT_SONNET_MODEL %s, got %s", p.Models.Sonnet, loaded.Env["ANTHROPIC_DEFAULT_SONNET_MODEL"])
	}
	if loaded.Env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] != p.Models.Haiku {
		t.Errorf("Expected ANTHROPIC_DEFAULT_HAIKU_MODEL %s, got %s", p.Models.Haiku, loaded.Env["ANTHROPIC_DEFAULT_HAIKU_MODEL"])
	}

	if loaded.Env["EXISTING_VAR"] != "value" {
		t.Error("Existing env var should be preserved")
	}
}

func TestGetCurrentProvider(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	settings := &Settings{
		Env: map[string]string{
			"ANTHROPIC_BASE_URL":   "https://api.test.com",
			"ANTHROPIC_AUTH_TOKEN": "test-token",
		},
		Permissions: Permissions{},
	}

	if err := WriteSettings(settings); err != nil {
		t.Fatalf("WriteSettings failed: %v", err)
	}

	baseURL, apiKey := GetCurrentProvider()
	if baseURL != "https://api.test.com" {
		t.Errorf("Expected baseURL 'https://api.test.com', got %s", baseURL)
	}
	if apiKey != "test-token" {
		t.Errorf("Expected apiKey 'test-token', got %s", apiKey)
	}
}

func TestValidateProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider *provider.Provider
		wantErr  bool
	}{
		{
			name: "valid provider",
			provider: &provider.Provider{
				BaseURL: "https://api.test.com",
				APIKey:  "test-key",
				Models: provider.ModelConfig{
					Opus:   "opus",
					Sonnet: "sonnet",
					Haiku:  "haiku",
				},
			},
			wantErr: false,
		},
		{
			name: "missing base_url",
			provider: &provider.Provider{
				BaseURL: "",
				APIKey:  "test-key",
				Models: provider.ModelConfig{
					Opus:   "opus",
					Sonnet: "sonnet",
					Haiku:  "haiku",
				},
			},
			wantErr: true,
		},
		{
			name: "missing api_key",
			provider: &provider.Provider{
				BaseURL: "https://api.test.com",
				APIKey:  "",
				Models: provider.ModelConfig{
					Opus:   "opus",
					Sonnet: "sonnet",
					Haiku:  "haiku",
				},
			},
			wantErr: true,
		},
		{
			name: "missing opus model",
			provider: &provider.Provider{
				BaseURL: "https://api.test.com",
				APIKey:  "test-key",
				Models: provider.ModelConfig{
					Opus:   "",
					Sonnet: "sonnet",
					Haiku:  "haiku",
				},
			},
			wantErr: true,
		},
		{
			name: "missing sonnet model",
			provider: &provider.Provider{
				BaseURL: "https://api.test.com",
				APIKey:  "test-key",
				Models: provider.ModelConfig{
					Opus:   "opus",
					Sonnet: "",
					Haiku:  "haiku",
				},
			},
			wantErr: true,
		},
		{
			name: "missing haiku model",
			provider: &provider.Provider{
				BaseURL: "https://api.test.com",
				APIKey:  "test-key",
				Models: provider.ModelConfig{
					Opus:   "opus",
					Sonnet: "sonnet",
					Haiku:  "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProvider(tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
