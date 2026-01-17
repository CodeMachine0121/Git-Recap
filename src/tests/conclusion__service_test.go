package tests

import (
	"commit-record/src/internal/domains"
	"commit-record/src/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConclusionByAi(t *testing.T) {

	proxy := NewMockOpenAiProxy()
	service := services.NewConclusionService(proxy)

	conclusionMessage := service.GetConclusion(domains.CommitRecord{
		ProjectName:   "project",
		CommitMessage: []string{"commit message1", "commit message2"},
	})

	assert.True(t, proxy.IsReceived)
	assert.Equal(t, "conclusion", conclusionMessage)
}
