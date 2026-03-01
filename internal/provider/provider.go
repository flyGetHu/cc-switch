package provider

// ModelConfig 与 Claude Code settings.json 中的 env 字段命名保持一致
type ModelConfig struct {
	DefaultOpus   string `yaml:"default_opus"`   // ANTHROPIC_DEFAULT_OPUS_MODEL
	DefaultSonnet string `yaml:"default_sonnet"` // ANTHROPIC_DEFAULT_SONNET_MODEL
	DefaultHaiku  string `yaml:"default_haiku"`  // ANTHROPIC_DEFAULT_HAIKU_MODEL
	SmallFast     string `yaml:"small_fast"`     // ANTHROPIC_SMALL_FAST_MODEL
	DefaultModel  string `yaml:"default"`        // ANTHROPIC_MODEL
}

type Provider struct {
	Name    string      `yaml:"name"`
	BaseURL string      `yaml:"base_url"`
	Models  ModelConfig `yaml:"models"`
	APIKey  string      `yaml:"api_key"`
}

func GetPresets() map[string]Provider {
	return map[string]Provider{
		"zhipu": {
			Name:    "智普 GLM",
			BaseURL: "https://open.bigmodel.cn/api/anthropic",
			Models: ModelConfig{
				DefaultOpus:   "glm-5",
				DefaultSonnet: "glm-5",
				DefaultHaiku:  "glm-5",
				SmallFast:     "glm-5",
				DefaultModel:  "glm-5",
			},
		},
		"minimax": {
			Name:    "MiniMax",
			BaseURL: "https://api.minimaxi.com/anthropic",
			Models: ModelConfig{
				DefaultOpus:   "MiniMax-M2.5",
				DefaultSonnet: "MiniMax-M2.5",
				DefaultHaiku:  "MiniMax-M2.5",
				SmallFast:     "MiniMax-M2.5",
				DefaultModel:  "MiniMax-M2.5",
			},
		},
	}
}
