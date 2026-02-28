package cmd

import (
	"fmt"
	"os"

	"cc-switch/internal/tui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cc-switch",
	Short: "Claude Code 模型服务商切换工具",
	Long:  "cc-switch 是一个用于切换 Claude Code 模型服务商的 CLI 工具",
	Run: func(cmd *cobra.Command, args []string) {
		choice, err := tui.RunSelector()
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}

		if choice == "" {
			fmt.Println("已取消")
			return
		}

		if err := switchProvider(choice); err != nil {
			fmt.Fprintf(os.Stderr, "切换失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("已切换到 %s\n", choice)
		fmt.Println("请重启 Claude Code 使配置生效")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(backupCmd)
}
