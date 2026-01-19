package tests

import (
	"commit-record/src/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDailyWorkConclusion(t *testing.T) {

	gitHandler := NewMockGitHandler()
	openAiProxy := NewMockOpenAiProxy()
	persistenceRepo := NewMockPersistenceRepository()

	service := services.NewDailyService(
		services.NewGitCommitService(gitHandler),
		services.NewConclusionService(openAiProxy),
		persistenceRepo)

	service.DoDailyWorkConclusion("any-author", "path/to/project")

	assert.True(t, gitHandler.IsReceived)
	assert.True(t, openAiProxy.IsReceived)
}

func TestGetDailyWorkConclusionBatch(t *testing.T) {

	gitHandler := NewMockGitHandler()
	openAiProxy := NewMockOpenAiProxy()
	persistenceRepo := NewMockPersistenceRepository()

	service := services.NewDailyService(
		services.NewGitCommitService(gitHandler),
		services.NewConclusionService(openAiProxy),
		persistenceRepo)

	service.DoDailyWorkConclusionBatch("any-author", []string{"path/to/project", "path/to/project2"})

	assert.True(t, gitHandler.IsReceived)
	assert.True(t, openAiProxy.IsReceived)
}
