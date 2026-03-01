package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"cc-switch/internal/claude"
	"cc-switch/internal/config"
)

var getBackupsDir = config.GetBackupsDir
var getClaudeSettingsPath = claude.GetSettingsPath

func CreateBackup() (string, error) {
	data, err := os.ReadFile(getClaudeSettingsPath())
	if err != nil {
		return "", err
	}

	backupsDir, err := getBackupsDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		return "", err
	}

	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("settings-%s.json", timestamp)
	backupPath := filepath.Join(backupsDir, filename)

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return "", err
	}

	return backupPath, nil
}

func ListBackups() ([]string, error) {
	backupsDir, err := getBackupsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(backupsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			backups = append(backups, filepath.Join(backupsDir, entry.Name()))
		}
	}

	// 按文件名（时间戳）升序排序，最新的在前面
	sort.Slice(backups, func(i, j int) bool {
		return backups[i] > backups[j]
	})

	return backups, nil
}

func RestoreBackup(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	return os.WriteFile(getClaudeSettingsPath(), data, 0600)
}

func GetLatestBackup() (string, error) {
	backups, err := ListBackups()
	if err != nil {
		return "", err
	}

	if len(backups) == 0 {
		return "", fmt.Errorf("没有可用的备份")
	}

	return backups[0], nil
}
