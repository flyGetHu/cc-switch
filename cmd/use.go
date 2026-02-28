package cmd

import (
	"fmt"
	"os"

	"cc-switch/internal/backup"
	"cc-switch/internal/claude"
	"cc-switch/internal/config"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <provider>",
	Short: "切换到指定服务商",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if err := switchProvider(name); err != nil {
			fmt.Fprintf(os.Stderr, "切换失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("已切换到 %s\n", name)
		fmt.Println("请重启 Claude Code 使配置生效")
	},
}

func switchProvider(name string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	p, ok := cfg.Providers[name]
	if !ok {
		return fmt.Errorf("服务商 '%s' 不存在", name)
	}

	if p.APIKey == "" {
		return fmt.Errorf("服务商 '%s' 未配置 API Key，请先使用 'cc-switch edit %s' 配置", name, name)
	}

	if _, err := backup.CreateBackup(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 创建备份失败: %v\n", err)
	}

	if err := claude.ApplyProvider(&p); err != nil {
		return fmt.Errorf("应用配置失败: %w", err)
	}

	if err := config.SetCurrent(name); err != nil {
		return fmt.Errorf("保存当前服务商失败: %w", err)
	}

	return nil
}
