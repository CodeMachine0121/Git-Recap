package tests

import (
	"commit-record/src/internal/services"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetDailyCommitMessage(t *testing.T) {
	handler := NewMockGitHandler()
	service := services.NewGitCommitService(handler)

	messages := service.GetDailyCommitMessages("path/to/project", "")

	assert.True(t, handler.IsReceived)
	assert.True(t, len(messages.CommitMessage) == 2)
}
