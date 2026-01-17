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

	prompt := fmt.Sprintf(`
Based on the following commit messages for the project "%s", please provide a concise summary of the work done today. Focus on the key achievements and changes made.

Commit messages:
%s

Please provide a brief conclusion (2-3 sentences) summarizing the main accomplishments ( using trandtional Chinese ).`,
		record.ProjectName,
		commitList.String())

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT5Mini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxCompletionTokens: 200,
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
