package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"cc-switch/internal/backup"
	"cc-switch/internal/config"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "备份管理",
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有备份",
	Run: func(cmd *cobra.Command, args []string) {
		backups, err := backup.ListBackups()
		if err != nil {
			fmt.Fprintf(os.Stderr, "获取备份列表失败: %v\n", err)
			return
		}

		if len(backups) == 0 {
			fmt.Println("没有备份")
			return
		}

		fmt.Println("备份列表:")
		for i, b := range backups {
			filename := filepath.Base(b)
			fmt.Printf("  %d. %s\n", i+1, filename)
		}
	},
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore [file]",
	Short: "恢复备份",
	Run: func(cmd *cobra.Command, args []string) {
		var backupPath string

		if len(args) > 0 {
			backupsDir := config.GetBackupsDir()
			backupPath = filepath.Join(backupsDir, args[0])
		} else {
			var err error
			backupPath, err = backup.GetLatestBackup()
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取最新备份失败: %v\n", err)
				return
			}
		}

		fmt.Printf("即将恢复备份: %s\n", filepath.Base(backupPath))
		fmt.Print("确认恢复? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("已取消")
			return
		}

		if err := backup.RestoreBackup(backupPath); err != nil {
			fmt.Fprintf(os.Stderr, "恢复失败: %v\n", err)
			return
		}

		fmt.Println("已恢复备份")
		fmt.Println("请重启 Claude Code 使配置生效")
	},
}

var backupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "创建备份",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := backup.CreateBackup()
		if err != nil {
			fmt.Fprintf(os.Stderr, "创建备份失败: %v\n", err)
			return
		}

		fmt.Printf("已创建备份: %s\n", filepath.Base(path))
	},
}

func init() {
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupCreateCmd)
}
