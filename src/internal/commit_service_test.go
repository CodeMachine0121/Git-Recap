package internal

import (
	"commit-record/src/internal/services"

	"github.com/stretchr/testify/assert"

	"testing"
)

type MockGitHandler struct {
	IsReceived bool
}

func (m *MockGitHandler) GetCommitMessages(string) []string {
	m.IsReceived = true
	return []string{"commit message1", "commit message2"}
}
func NewMockGitHandler() *MockGitHandler {
	return &MockGitHandler{}
}

func TestGetDailyCommitMessage(t *testing.T) {
	handler := NewMockGitHandler()
	service := services.NewDailyService(handler)

	messages := service.GetDailyCommitMessages("path/to/project")

	assert.True(t, handler.IsReceived)
	assert.True(t, len(messages.CommitMessage) == 2)
}
