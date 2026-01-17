package proxies

import (
	"commit-record/src/internal/domains"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type IOpenAiProxy interface {
	GetConclusion(record domains.CommitRecord) string
}

type OpenAiProxy struct {
	client *openai.Client
}

func NewOpenAiProxy() *OpenAiProxy {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)
	return &OpenAiProxy{
		client: client,
	}
}

func (p *OpenAiProxy) GetConclusion(record domains.CommitRecord) string {
	// Build the prompt with commit messages
	var commitList strings.Builder
	for i, msg := range record.CommitMessage {
		commitList.WriteString(fmt.Sprintf("%d. %s\n", i+1, msg))
	}

	prompt := fmt.Sprintf(`專案「%s」的今日 commit 訊息：

%s

請用繁體中文總結今日主要成果（2-3句話）。`,
		record.ProjectName,
		commitList.String())

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxCompletionTokens: 150,
			Temperature:         0.3,
		},
	)

	if err != nil {
		return fmt.Sprintf("Error generating conclusion: %v", err)
	}

	if len(resp.Choices) > 0 {
		return strings.TrimSpace(resp.Choices[0].Message.Content)
	}

	return "Unable to generate conclusion"
}
