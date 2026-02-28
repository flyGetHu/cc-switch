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

func TestReadSettingsRaw_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	s, err := ReadSettingsRaw()
	if err != nil {
		t.Errorf("ReadSettingsRaw failed: %v", err)
	}

	if s == nil {
		t.Error("Settings should not be nil")
	}
}

func TestReadAndWriteSettingsRaw_PreserveUnknownFields(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	// 模拟包含未知字段的 settings.json
	original := map[string]interface{}{
		"env": map[string]interface{}{
			"TEST_VAR": "test_value",
		},
		"permissions": map[string]interface{}{
			"allow": []string{"read"},
			"deny":  []string{"write"},
		},
		"apiProvider":   "anthropic",
		"primaryApiKey": "sk-ant-key",
		"unknownField": 123,
	}

	if err := WriteSettingsRaw(original); err != nil {
		t.Fatalf("WriteSettingsRaw failed: %v", err)
	}

	loaded, err := ReadSettingsRaw()
	if err != nil {
		t.Fatalf("ReadSettingsRaw failed: %v", err)
	}

	// 验证 env 字段
	env, ok := loaded["env"].(map[string]interface{})
	if !ok {
		t.Fatal("env should be a map")
	}
	if env["TEST_VAR"] != "test_value" {
		t.Errorf("Expected TEST_VAR 'test_value', got %v", env["TEST_VAR"])
	}

	// 验证 permissions 字段
	perms, ok := loaded["permissions"].(map[string]interface{})
	if !ok {
		t.Fatal("permissions should be a map")
	}
	allow, ok := perms["allow"].([]interface{})
	if !ok || len(allow) != 1 || allow[0] != "read" {
		t.Errorf("Permissions.Allow not preserved correctly")
	}

	// 验证未知字段被保留
	if loaded["apiProvider"] != "anthropic" {
		t.Errorf("apiProvider not preserved: got %v", loaded["apiProvider"])
	}
	if loaded["primaryApiKey"] != "sk-ant-key" {
		t.Errorf("primaryApiKey not preserved: got %v", loaded["primaryApiKey"])
	}
	if loaded["unknownField"] != float64(123) {
		t.Errorf("unknownField not preserved: got %v", loaded["unknownField"])
	}
}

func TestApplyProvider_PreserveOtherFields(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	// 初始配置包含其他字段
	original := map[string]interface{}{
		"env": map[string]interface{}{
			"EXISTING_VAR":     "value",
			"ANOTHER_VAR":     "keep_this",
			"ANTHROPIC_BASE_URL": "https://old-api.com",
		},
		"permissions": map[string]interface{}{
			"allow": []string{"read"},
		},
		"apiProvider":   "anthropic",
		"primaryApiKey": "sk-ant-key",
	}

	if err := WriteSettingsRaw(original); err != nil {
		t.Fatalf("WriteSettingsRaw failed: %v", err)
	}

	p := &provider.Provider{
		Name:    "Test",
		BaseURL: "https://api.test.com/anthropic",
		Models: provider.ModelConfig{
			DefaultOpus:   "test-opus",
			DefaultSonnet: "test-sonnet",
			DefaultHaiku:  "test-haiku",
			SmallFast:     "test-small-fast",
			DefaultModel:  "test-default",
		},
		APIKey: "test-api-key",
	}

	if err := ApplyProvider(p); err != nil {
		t.Fatalf("ApplyProvider failed: %v", err)
	}

	loaded, err := ReadSettingsRaw()
	if err != nil {
		t.Fatalf("ReadSettingsRaw failed: %v", err)
	}

	// 验证服务商相关字段被更新
	env, ok := loaded["env"].(map[string]interface{})
	if !ok {
		t.Fatal("env should be a map")
	}

	if env["ANTHROPIC_BASE_URL"] != p.BaseURL {
		t.Errorf("Expected ANTHROPIC_BASE_URL %s, got %v", p.BaseURL, env["ANTHROPIC_BASE_URL"])
	}
	if env["ANTHROPIC_AUTH_TOKEN"] != p.APIKey {
		t.Errorf("Expected ANTHROPIC_AUTH_TOKEN %s, got %v", p.APIKey, env["ANTHROPIC_AUTH_TOKEN"])
	}
	if env["ANTHROPIC_DEFAULT_OPUS_MODEL"] != p.Models.DefaultOpus {
		t.Errorf("Expected ANTHROPIC_DEFAULT_OPUS_MODEL %s, got %v", p.Models.DefaultOpus, env["ANTHROPIC_DEFAULT_OPUS_MODEL"])
	}
	if env["ANTHROPIC_DEFAULT_SONNET_MODEL"] != p.Models.DefaultSonnet {
		t.Errorf("Expected ANTHROPIC_DEFAULT_SONNET_MODEL %s, got %v", p.Models.DefaultSonnet, env["ANTHROPIC_DEFAULT_SONNET_MODEL"])
	}
	if env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] != p.Models.DefaultHaiku {
		t.Errorf("Expected ANTHROPIC_DEFAULT_HAIKU_MODEL %s, got %v", p.Models.DefaultHaiku, env["ANTHROPIC_DEFAULT_HAIKU_MODEL"])
	}
	if env["ANTHROPIC_SMALL_FAST_MODEL"] != p.Models.SmallFast {
		t.Errorf("Expected ANTHROPIC_SMALL_FAST_MODEL %s, got %v", p.Models.SmallFast, env["ANTHROPIC_SMALL_FAST_MODEL"])
	}
	if env["ANTHROPIC_MODEL"] != p.Models.DefaultModel {
		t.Errorf("Expected ANTHROPIC_MODEL %s, got %v", p.Models.DefaultModel, env["ANTHROPIC_MODEL"])
	}

	// 验证其他非服务商字段被保留
	if env["EXISTING_VAR"] != "value" {
		t.Error("Existing env var should be preserved")
	}
	if env["ANOTHER_VAR"] != "keep_this" {
		t.Error("ANOTHER_VAR should be preserved")
	}

	// 验证其他顶层字段被保留
	if loaded["apiProvider"] != "anthropic" {
		t.Errorf("apiProvider not preserved: got %v", loaded["apiProvider"])
	}
	if loaded["primaryApiKey"] != "sk-ant-key" {
		t.Errorf("primaryApiKey not preserved: got %v", loaded["primaryApiKey"])
	}
}

func TestGetCurrentProvider(t *testing.T) {
	tempDir := t.TempDir()
	tempPath := filepath.Join(tempDir, "settings.json")

	SetSettingsPath(tempPath)
	defer SetSettingsPath("")

	settings := map[string]interface{}{
		"env": map[string]interface{}{
			"ANTHROPIC_BASE_URL":   "https://api.test.com",
			"ANTHROPIC_AUTH_TOKEN": "test-token",
		},
		"permissions": map[string]interface{}{},
	}

	if err := WriteSettingsRaw(settings); err != nil {
		t.Fatalf("WriteSettingsRaw failed: %v", err)
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
					DefaultOpus:   "opus",
					DefaultSonnet: "sonnet",
					DefaultHaiku:  "haiku",
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
					DefaultOpus:   "opus",
					DefaultSonnet: "sonnet",
					DefaultHaiku:  "haiku",
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
					DefaultOpus:   "opus",
					DefaultSonnet: "sonnet",
					DefaultHaiku:  "haiku",
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
					DefaultOpus:   "",
					DefaultSonnet: "sonnet",
					DefaultHaiku:  "haiku",
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
					DefaultOpus:   "opus",
					DefaultSonnet: "",
					DefaultHaiku:  "haiku",
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
					DefaultOpus:   "opus",
					DefaultSonnet: "sonnet",
					DefaultHaiku:  "",
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
