package git

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type IGitHandler interface {
	GetCommitMessages(projectPath string) []string
}

type GitHandler struct{}

func (*GitHandler) GetCommitMessages(projectPath string) []string {
	// Get today's date at midnight for filtering commits
	today := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	cmd := exec.Command("git", "--no-pager",
		"log",
		"--since="+today,
		"--format=%s")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		// 這裡會印出真正的 Git 錯誤，例如 "fatal: not a git repository"
		fmt.Printf("Git Error Message: %s\n", string(output))
		fmt.Printf("Error Code: %v\n", err)
		return nil
	}

	// Split output by newlines and filter empty strings
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	// Filter out empty lines
	var messages []string
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			messages = append(messages, trimmed)
		}
	}

	return messages
}
