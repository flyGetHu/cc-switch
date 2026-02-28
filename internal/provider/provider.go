package provider

type ModelConfig struct {
	Opus   string `yaml:"opus"`
	Sonnet string `yaml:"sonnet"`
	Haiku  string `yaml:"haiku"`
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
				Opus:   "glm-5",
				Sonnet: "glm-5",
				Haiku:  "glm-5",
			},
		},
		"minimax": {
			Name:    "MiniMax",
			BaseURL: "https://api.minimaxi.com/anthropic",
			Models: ModelConfig{
				Opus:   "MiniMax-M2.5",
				Sonnet: "MiniMax-M2.5",
				Haiku:  "MiniMax-M2.5",
			},
		},
	}
}
