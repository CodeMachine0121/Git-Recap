package tests

import (
	"commit-record/src/internal/domains"
)

type MockOpenAiProxy struct{ IsReceived bool }

func (m *MockOpenAiProxy) GetConclusion(domains.CommitRecord) string {
	m.IsReceived = true
	return "conclusion"
}

func (m *MockOpenAiProxy) GetBatchConclusion([]domains.CommitRecord) map[string]string {
	m.IsReceived = true
	return map[string]string{"project1": "conclusion1", "project2": "conclusion2"}
}

func NewMockOpenAiProxy() *MockOpenAiProxy {
	return &MockOpenAiProxy{}
}

type MockGitHandler struct {
	IsReceived bool
}

func (m *MockGitHandler) GetCommitMessages(string, string) []string {
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
