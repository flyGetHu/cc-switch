package cmd

import (
	"fmt"
	"os"

	"cc-switch/internal/config"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <provider>",
	Short: "编辑服务商配置",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
			return
		}

		p, ok := cfg.Providers[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "服务商 '%s' 不存在\n", name)
			return
		}

		changed := false

		if apiKey, _ := cmd.Flags().GetString("api-key"); apiKey != "" {
			p.APIKey = apiKey
			changed = true
		}

		if baseURL, _ := cmd.Flags().GetString("base-url"); baseURL != "" {
			p.BaseURL = baseURL
			changed = true
		}

		if opus, _ := cmd.Flags().GetString("opus"); opus != "" {
			p.Models.DefaultOpus = opus
			changed = true
		}

		if sonnet, _ := cmd.Flags().GetString("sonnet"); sonnet != "" {
			p.Models.DefaultSonnet = sonnet
			changed = true
		}

		if haiku, _ := cmd.Flags().GetString("haiku"); haiku != "" {
			p.Models.DefaultHaiku = haiku
			changed = true
		}

		if !changed {
			fmt.Println("请指定要修改的选项:")
			fmt.Println("  --api-key    API Key")
			fmt.Println("  --base-url   Base URL")
			fmt.Println("  --opus       Opus 模型 (default_opus)")
			fmt.Println("  --sonnet     Sonnet 模型 (default_sonnet)")
			fmt.Println("  --haiku      Haiku 模型 (default_haiku)")
			return
		}

		if err := config.AddProvider(name, p); err != nil {
			fmt.Fprintf(os.Stderr, "保存失败: %v\n", err)
			return
		}

		fmt.Printf("已更新服务商: %s\n", name)
	},
}

func init() {
	editCmd.Flags().StringP("api-key", "k", "", "API Key")
	editCmd.Flags().StringP("base-url", "u", "", "Base URL")
	editCmd.Flags().StringP("opus", "", "", "Opus 模型")
	editCmd.Flags().StringP("sonnet", "", "", "Sonnet 模型")
	editCmd.Flags().StringP("haiku", "", "", "Haiku 模型")
}
