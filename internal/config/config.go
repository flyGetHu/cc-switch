package config

import (
	"fmt"
	"os"
	"path/filepath"

	"cc-switch/internal/provider"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers  map[string]provider.Provider `yaml:"providers"`
	Current    string                       `yaml:"current"`
	BackupsDir string                       `yaml:"backups_dir"`
}

var configPath string

func GetConfigPath() string {
	if configPath != "" {
		return configPath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "cc-switch", "config.yaml")
}

func SetConfigPath(path string) {
	configPath = path
}

func GetBackupsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "cc-switch", "backups")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(GetConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return initDefaultConfig()
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func initDefaultConfig() (*Config, error) {
	presets := provider.GetPresets()
	cfg := &Config{
		Providers:  presets,
		Current:    "zhipu",
		BackupsDir: GetBackupsDir(),
	}

	if err := Save(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func Save(c *Config) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	cp := GetConfigPath()
	configDir := filepath.Dir(cp)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(cp, data, 0644)
}

func GetCurrent() (*Config, *provider.Provider, error) {
	cfg, err := Load()
	if err != nil {
		return nil, nil, err
	}

	p, ok := cfg.Providers[cfg.Current]
	if !ok {
		return cfg, nil, nil
	}

	return cfg, &p, nil
}

func AddProvider(name string, p provider.Provider) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if cfg.Providers == nil {
		cfg.Providers = make(map[string]provider.Provider)
	}
	cfg.Providers[name] = p

	return Save(cfg)
}

func RemoveProvider(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	delete(cfg.Providers, name)

	if cfg.Current == name {
		for k := range cfg.Providers {
			cfg.Current = k
			break
		}
	}

	return Save(cfg)
}

func SetCurrent(name string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	if _, ok := cfg.Providers[name]; !ok {
		return ErrProviderNotFound
	}

	cfg.Current = name
	return Save(cfg)
}

func UpdateProviderAPIKey(name, apiKey string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	p, ok := cfg.Providers[name]
	if !ok {
		return ErrProviderNotFound
	}

	p.APIKey = apiKey
	cfg.Providers[name] = p

	return Save(cfg)
}

var ErrProviderNotFound = fmt.Errorf("provider not found")
