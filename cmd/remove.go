package cmd

import (
	"fmt"
	"os"

	"cc-switch/internal/config"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <provider>",
	Short: "删除服务商",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
			return
		}

		if _, ok := cfg.Providers[name]; !ok {
			fmt.Fprintf(os.Stderr, "服务商 '%s' 不存在\n", name)
			return
		}

		if name == "zhipu" || name == "minimax" {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				fmt.Println("警告: 这是预设服务商，删除后可以使用 'cc-switch add' 重新添加")
				fmt.Print("确认删除? (y/N): ")
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("已取消")
					return
				}
			}
		}

		if err := config.RemoveProvider(name); err != nil {
			fmt.Fprintf(os.Stderr, "删除失败: %v\n", err)
			return
		}

		fmt.Printf("已删除服务商: %s\n", name)
	},
}

func init() {
	removeCmd.Flags().BoolP("force", "f", false, "强制删除预设服务商")
}
