package cmd

import (
	"fmt"

	"cc-switch/internal/config"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "显示当前服务商",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, p, err := config.GetCurrent()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "加载配置失败: %v\n", err)
			return
		}

		if p == nil {
			fmt.Println("未配置当前服务商")
			return
		}

		fmt.Printf("当前服务商: %s (%s)\n", cfg.Current, p.Name)
		fmt.Printf("Base URL: %s\n", p.BaseURL)
		fmt.Printf("Opus: %s, Sonnet: %s, Haiku: %s (default: %s)\n",
			p.Models.DefaultOpus, p.Models.DefaultSonnet, p.Models.DefaultHaiku, p.Models.DefaultModel)

		if p.APIKey == "" {
			fmt.Println("API Key: 未配置")
		} else {
			maskedKey := maskAPIKey(p.APIKey)
			fmt.Printf("API Key: %s\n", maskedKey)
		}
	},
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
