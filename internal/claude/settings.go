package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cc-switch/internal/provider"
)

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

// ReadSettingsRaw 读取完整的 settings.json 为 map，保留所有字段
func ReadSettingsRaw() (map[string]interface{}, error) {
	data, err := os.ReadFile(GetSettingsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]interface{}), nil
		}
		return nil, err
	}

	var s map[string]interface{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	if s == nil {
		s = make(map[string]interface{})
	}

	return s, nil
}

// WriteSettingsRaw 写入完整的 settings map
func WriteSettingsRaw(s map[string]interface{}) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetSettingsPath(), data, 0600)
}

// getEnvMap 从 settings 中获取 env map，如果不存在则创建
func getEnvMap(s map[string]interface{}) map[string]interface{} {
	if s["env"] == nil {
		s["env"] = make(map[string]interface{})
	}
	env, ok := s["env"].(map[string]interface{})
	if !ok {
		// 如果 env 不是 map，重新创建
		env = make(map[string]interface{})
		s["env"] = env
	}
	return env
}

// ApplyProvider 应用服务商配置，只修改 env 中的服务商相关字段，保留其他所有配置
func ApplyProvider(p *provider.Provider) error {
	s, err := ReadSettingsRaw()
	if err != nil {
		return err
	}

	env := getEnvMap(s)

	// 只更新服务商相关的环境变量（与 Claude Code settings.json 命名一致）
	env["ANTHROPIC_BASE_URL"] = p.BaseURL
	env["ANTHROPIC_AUTH_TOKEN"] = p.APIKey
	env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = p.Models.DefaultOpus
	env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = p.Models.DefaultSonnet
	env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = p.Models.DefaultHaiku
	env["ANTHROPIC_SMALL_FAST_MODEL"] = p.Models.SmallFast
	env["ANTHROPIC_MODEL"] = p.Models.DefaultModel

	return WriteSettingsRaw(s)
}

func GetCurrentProvider() (string, string) {
	s, err := ReadSettingsRaw()
	if err != nil {
		return "", ""
	}

	env := getEnvMap(s)

	baseURL, _ := env["ANTHROPIC_BASE_URL"].(string)
	apiKey, _ := env["ANTHROPIC_AUTH_TOKEN"].(string)

	return baseURL, apiKey
}

func ValidateProvider(p *provider.Provider) error {
	if p.BaseURL == "" {
		return fmt.Errorf("base_url 不能为空")
	}
	if p.APIKey == "" {
		return fmt.Errorf("api_key 不能为空")
	}
	if p.Models.DefaultOpus == "" || p.Models.DefaultSonnet == "" || p.Models.DefaultHaiku == "" {
		return fmt.Errorf("模型配置不完整")
	}
	// SmallFast 和 General 可以为空，使用默认值
	return nil
}
