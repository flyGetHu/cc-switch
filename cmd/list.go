package cmd

import (
	"fmt"

	"cc-switch/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有服务商",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "加载配置失败: %v\n", err)
			return
		}

		fmt.Println("可用服务商:")
		fmt.Println()

		for key, p := range cfg.Providers {
			current := ""
			if key == cfg.Current {
				current = " (当前)"
			}

			apiKeyStatus := "未配置"
			if p.APIKey != "" {
				apiKeyStatus = "已配置"
			}

			fmt.Printf("  %-10s %-15s %s%s\n", key, p.Name, apiKeyStatus, current)
			fmt.Printf("             Base URL: %s\n", p.BaseURL)
			fmt.Printf("             Opus: %s, Sonnet: %s, Haiku: %s\n",
				p.Models.Opus, p.Models.Sonnet, p.Models.Haiku)
			fmt.Println()
		}
	},
}
