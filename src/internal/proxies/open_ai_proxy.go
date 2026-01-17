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
	GetBatchConclusion(records []domains.CommitRecord) map[string]string
}

type OpenAiProxy struct {
	client *openai.Client
	parser *MultiProjectResponseParser
}

func NewOpenAiProxy() *OpenAiProxy {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)
	return &OpenAiProxy{
		client: client,
		parser: NewMultiProjectResponseParser(),
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

func (p *OpenAiProxy) GetBatchConclusion(records []domains.CommitRecord) map[string]string {

	// 合併所有專案成一個請求
	var promptBuilder strings.Builder
	promptBuilder.WriteString("以下是多個專案的今日 commit 訊息，請分別總結每個專案的主要成果（每個專案 2-3 句話）：\n\n")

	for _, record := range records {
		promptBuilder.WriteString(fmt.Sprintf("## 專案「%s」\n", record.ProjectName))
		for i, msg := range record.CommitMessage {
			promptBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, msg))
		}
		promptBuilder.WriteString("\n")
	}

	promptBuilder.WriteString("請以以下格式回覆，每個專案一段：\n")
	promptBuilder.WriteString("【專案名稱】\n總結內容\n\n")

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: promptBuilder.String(),
				},
			},
			MaxCompletionTokens: 500,
			Temperature:         0.3,
		},
	)

	if err != nil {
		// 錯誤時返回錯誤訊息
		result := make(map[string]string)
		for _, record := range records {
			result[record.ProjectName] = fmt.Sprintf("Error: %v", err)
		}
		return result
	}

	// 解析回覆
	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
	}

	return p.parseMultiProjectResponse(content, records)
}

func (p *OpenAiProxy) parseMultiProjectResponse(content string, records []domains.CommitRecord) map[string]string {
	return p.parser.Parse(content, records)
}
