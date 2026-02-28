package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cc-switch/internal/provider"
)

type Settings struct {
	Env         map[string]string `json:"env"`
	Permissions Permissions       `json:"permissions"`
}

type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

var settingsPath string

func GetSettingsPath() string {
	if settingsPath != "" {
		return settingsPath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", "settings.json")
}

func SetSettingsPath(path string) {
	settingsPath = path
}

func ReadSettings() (*Settings, error) {
	data, err := os.ReadFile(GetSettingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return &Settings{
				Env:         make(map[string]string),
				Permissions: Permissions{Allow: []string{}, Deny: []string{}},
			}, nil
		}
		return nil, err
	}

	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	if s.Env == nil {
		s.Env = make(map[string]string)
	}

	return &s, nil
}

func WriteSettings(s *Settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetSettingsPath(), data, 0644)
}

func ApplyProvider(p *provider.Provider) error {
	s, err := ReadSettings()
	if err != nil {
		return err
	}

	s.Env["ANTHROPIC_BASE_URL"] = p.BaseURL
	s.Env["ANTHROPIC_AUTH_TOKEN"] = p.APIKey
	s.Env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = p.Models.Opus
	s.Env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = p.Models.Sonnet
	s.Env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = p.Models.Haiku

	return WriteSettings(s)
}

func GetCurrentProvider() (string, string) {
	s, err := ReadSettings()
	if err != nil {
		return "", ""
	}

	baseURL := s.Env["ANTHROPIC_BASE_URL"]
	apiKey := s.Env["ANTHROPIC_AUTH_TOKEN"]

	return baseURL, apiKey
}

func ValidateProvider(p *provider.Provider) error {
	if p.BaseURL == "" {
		return fmt.Errorf("base_url 不能为空")
	}
	if p.APIKey == "" {
		return fmt.Errorf("api_key 不能为空")
	}
	if p.Models.Opus == "" || p.Models.Sonnet == "" || p.Models.Haiku == "" {
		return fmt.Errorf("模型配置不完整")
	}
	return nil
}
