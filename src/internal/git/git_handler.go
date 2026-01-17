package git

import (
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
	today := time.Now().Format("2006-01-02")

	// Execute git log command to get today's commit messages
	// --since=today --format=%s returns only commit subjects (one per line)
	cmd := exec.Command("git", "-C", projectPath, "log", "--since="+today, "--format=%s")

	output, err := cmd.Output()
	if err != nil {
		// Return empty slice if git command fails (not a git repo, invalid path, etc.)
		return []string{}
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
