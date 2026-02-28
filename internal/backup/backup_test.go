package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestBackupEnv(t *testing.T) (string, func()) {
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, "cc-switch")
	backupsDir := filepath.Join(configDir, "backups")
	claudeDir := filepath.Join(tempDir, ".claude")

	if err := os.MkdirAll(backupsDir, 0755); err != nil {
		t.Fatalf("Failed to create backups dir: %v", err)
	}
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("Failed to create claude dir: %v", err)
	}

	settingsContent := `{"env": {"ANTHROPIC_AUTH_TOKEN": "test-key"}}`
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if err := os.WriteFile(settingsPath, []byte(settingsContent), 0644); err != nil {
		t.Fatalf("Failed to write settings: %v", err)
	}

	oldGetBackupsDir := getBackupsDir
	oldGetClaudeSettingsPath := getClaudeSettingsPath

	getBackupsDir = func() string { return backupsDir }
	getClaudeSettingsPath = func() string { return settingsPath }

	cleanup := func() {
		getBackupsDir = oldGetBackupsDir
		getClaudeSettingsPath = oldGetClaudeSettingsPath
	}

	return tempDir, cleanup
}

func TestCreateBackup(t *testing.T) {
	_, cleanup := setupTestBackupEnv(t)
	defer cleanup()

	backupPath, err := CreateBackup()
	if err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("Backup file not created at %s", backupPath)
	}

	content, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup: %v", err)
	}

	expected := `{"env": {"ANTHROPIC_AUTH_TOKEN": "test-key"}}`
	if string(content) != expected {
		t.Errorf("Expected content %s, got %s", expected, string(content))
	}
}

func TestListBackups(t *testing.T) {
	_, cleanup := setupTestBackupEnv(t)
	defer cleanup()

	backupsDir := getBackupsDir()

	for i := 0; i < 3; i++ {
		filename := filepath.Join(backupsDir, "settings-20060102-15040"+string(rune('0'+i))+".json")
		if err := os.WriteFile(filename, []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to create test backup: %v", err)
		}
	}

	backups, err := ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) != 3 {
		t.Errorf("Expected 3 backups, got %d", len(backups))
	}
}

func TestListBackups_Empty(t *testing.T) {
	tempDir := t.TempDir()
	emptyDir := filepath.Join(tempDir, "empty")
	os.MkdirAll(emptyDir, 0755)

	oldDir := getBackupsDir
	defer func() { getBackupsDir = oldDir }()
	getBackupsDir = func() string { return emptyDir }

	backups, err := ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}

	if len(backups) != 0 {
		t.Errorf("Expected 0 backups, got %d", len(backups))
	}
}

func TestRestoreBackup(t *testing.T) {
	_, cleanup := setupTestBackupEnv(t)
	defer cleanup()

	backupsDir := getBackupsDir()
	settingsPath := getClaudeSettingsPath()

	backupContent := `{"env": {"ANTHROPIC_AUTH_TOKEN": "restored-key"}}`
	backupPath := filepath.Join(backupsDir, "settings-20060102-150405.json")
	if err := os.WriteFile(backupPath, []byte(backupContent), 0644); err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	if err := RestoreBackup(backupPath); err != nil {
		t.Fatalf("RestoreBackup failed: %v", err)
	}

	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("Failed to read restored settings: %v", err)
	}

	if string(content) != backupContent {
		t.Errorf("Expected content %s, got %s", backupContent, string(content))
	}
}

func TestGetLatestBackup(t *testing.T) {
	_, cleanup := setupTestBackupEnv(t)
	defer cleanup()

	backupsDir := getBackupsDir()

	files := []string{
		"settings-20260101-100000.json",
		"settings-20260102-150000.json",
		"settings-20260102-160000.json",
	}

	for _, f := range files {
		if err := os.WriteFile(filepath.Join(backupsDir, f), []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to create test backup: %v", err)
		}
	}

	latest, err := GetLatestBackup()
	if err != nil {
		t.Fatalf("GetLatestBackup failed: %v", err)
	}

	expected := filepath.Join(backupsDir, "settings-20260102-160000.json")
	if latest != expected {
		t.Errorf("Expected %s, got %s", expected, latest)
	}
}

func TestGetLatestBackup_NoBackups(t *testing.T) {
	tempDir := t.TempDir()
	emptyDir := filepath.Join(tempDir, "empty")
	os.MkdirAll(emptyDir, 0755)

	oldDir := getBackupsDir
	defer func() { getBackupsDir = oldDir }()
	getBackupsDir = func() string { return emptyDir }

	_, err := GetLatestBackup()
	if err == nil {
		t.Error("Expected error when no backups exist")
	}
}
