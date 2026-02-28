package config

import (
	"os"
	"path/filepath"
	"testing"

	"cc-switch/internal/provider"

	"gopkg.in/yaml.v3"
)

func setupTestConfig(t *testing.T) string {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Providers: map[string]provider.Provider{
			"test": {
				Name:    "Test Provider",
				BaseURL: "https://api.test.com/anthropic",
				Models: provider.ModelConfig{
					Opus:   "test-opus",
					Sonnet: "test-sonnet",
					Haiku:  "test-haiku",
				},
				APIKey: "test-key",
			},
		},
		Current:    "test",
		BackupsDir: filepath.Join(tempDir, "backups"),
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	return configPath
}

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "cc-switch", "config.yaml")

	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestSetConfigPath(t *testing.T) {
	testPath := "/tmp/test/config.yaml"
	SetConfigPath(testPath)

	if GetConfigPath() != testPath {
		t.Errorf("Expected %s, got %s", testPath, GetConfigPath())
	}

	SetConfigPath("")
}

func TestGetBackupsDir(t *testing.T) {
	dir := GetBackupsDir()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "cc-switch", "backups")

	if dir != expected {
		t.Errorf("Expected %s, got %s", expected, dir)
	}
}

func TestLoad(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Current != "test" {
		t.Errorf("Expected current 'test', got %s", cfg.Current)
	}

	p, ok := cfg.Providers["test"]
	if !ok {
		t.Fatal("Expected provider 'test' to exist")
	}

	if p.Name != "Test Provider" {
		t.Errorf("Expected name 'Test Provider', got %s", p.Name)
	}
}

func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfg := &Config{
		Providers: map[string]provider.Provider{
			"save-test": {
				Name:    "Save Test",
				BaseURL: "https://api.save-test.com",
				Models: provider.ModelConfig{
					Opus:   "save-opus",
					Sonnet: "save-sonnet",
					Haiku:  "save-haiku",
				},
				APIKey: "save-key",
			},
		},
		Current:    "save-test",
		BackupsDir: filepath.Join(tempDir, "backups"),
	}

	SetConfigPath(configPath)
	defer SetConfigPath("")

	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	var loaded Config
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if loaded.Current != "save-test" {
		t.Errorf("Expected current 'save-test', got %s", loaded.Current)
	}
}

func TestAddProvider(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	newProvider := provider.Provider{
		Name:    "New Provider",
		BaseURL: "https://api.new.com",
		Models: provider.ModelConfig{
			Opus:   "new-opus",
			Sonnet: "new-sonnet",
			Haiku:  "new-haiku",
		},
		APIKey: "new-key",
	}

	if err := AddProvider("new", newProvider); err != nil {
		t.Fatalf("AddProvider failed: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	p, ok := cfg.Providers["new"]
	if !ok {
		t.Fatal("Expected provider 'new' to exist")
	}

	if p.Name != "New Provider" {
		t.Errorf("Expected name 'New Provider', got %s", p.Name)
	}
}

func TestRemoveProvider(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	cfg, _ := Load()
	cfg.Providers["to-remove"] = provider.Provider{
		Name:    "To Remove",
		BaseURL: "https://api.remove.com",
		Models:  provider.ModelConfig{Opus: "o", Sonnet: "s", Haiku: "h"},
	}
	Save(cfg)

	if err := RemoveProvider("to-remove"); err != nil {
		t.Fatalf("RemoveProvider failed: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if _, ok := cfg.Providers["to-remove"]; ok {
		t.Error("Expected provider 'to-remove' to be removed")
	}
}

func TestSetCurrent(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	cfg, _ := Load()
	cfg.Providers["another"] = provider.Provider{
		Name:    "Another",
		BaseURL: "https://api.another.com",
		Models:  provider.ModelConfig{Opus: "o", Sonnet: "s", Haiku: "h"},
	}
	Save(cfg)

	if err := SetCurrent("another"); err != nil {
		t.Fatalf("SetCurrent failed: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Current != "another" {
		t.Errorf("Expected current 'another', got %s", cfg.Current)
	}
}

func TestSetCurrent_NonExistent(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	err := SetCurrent("non-existent")
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

func TestUpdateProviderAPIKey(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	if err := UpdateProviderAPIKey("test", "new-api-key"); err != nil {
		t.Fatalf("UpdateProviderAPIKey failed: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	p := cfg.Providers["test"]
	if p.APIKey != "new-api-key" {
		t.Errorf("Expected APIKey 'new-api-key', got %s", p.APIKey)
	}
}

func TestUpdateProviderAPIKey_NonExistent(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	err := UpdateProviderAPIKey("non-existent", "key")
	if err != ErrProviderNotFound {
		t.Errorf("Expected ErrProviderNotFound, got %v", err)
	}
}

func TestGetCurrent(t *testing.T) {
	tempPath := setupTestConfig(t)

	SetConfigPath(tempPath)
	defer SetConfigPath("")

	cfg, p, err := GetCurrent()
	if err != nil {
		t.Fatalf("GetCurrent failed: %v", err)
	}

	if cfg.Current != "test" {
		t.Errorf("Expected current 'test', got %s", cfg.Current)
	}

	if p == nil {
		t.Fatal("Expected provider to not be nil")
	}

	if p.Name != "Test Provider" {
		t.Errorf("Expected name 'Test Provider', got %s", p.Name)
	}
}
