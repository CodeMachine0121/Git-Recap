package tests

import (
	"commit-record/src/internal/domains"
)

type MockOpenAiProxy struct{ IsReceived bool }

func (m *MockOpenAiProxy) GetConclusion(domains.CommitRecord) string {
	m.IsReceived = true
	return "conclusion"
}
func NewMockOpenAiProxy() *MockOpenAiProxy {
	return &MockOpenAiProxy{}
}

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

type MockPersistenceRepository struct {
	IsReceived bool
}

func (m *MockPersistenceRepository) Save(domains.DailyWorkConclusionRecord) error {
	m.IsReceived = true
	return nil
}

func NewMockPersistenceRepository() *MockPersistenceRepository {
	return &MockPersistenceRepository{}
}
