package repositories

import (
	"commit-record/src/internal/domains"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type IPersistenceRepository interface {
	Save(record domains.DailyWorkConclusionRecord) error
}

type PersistenceRepository struct{}

func (*PersistenceRepository) Save(record domains.DailyWorkConclusionRecord) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create directory path with current date (YYYY-MM-DD format)
	currentDate := time.Now().Format("2006-01-02")
	dirPath := filepath.Join(homeDir, "work_conclusion", currentDate)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	// Create file name based on project name
	fileName := fmt.Sprintf("%s.md", record.ProjectName)
	filePath := filepath.Join(dirPath, fileName)

	// Prepare markdown content with structured format
	content := fmt.Sprintf("# 每日工作總結 - %s\n\n", currentDate)
	content += fmt.Sprintf("## 專案：%s\n\n", record.ProjectName)
	content += record.Conclusion + "\n"

	// Write to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	fmt.Printf("✓ 已儲存工作總結：%s\n", filePath)
	return nil
}
