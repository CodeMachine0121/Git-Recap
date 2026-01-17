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

	prompt := fmt.Sprintf(`請根據專案「%s」的今日 commit 訊息，生成一份結構化的工作總結：

%s

請以以下格式用繁體中文輸出（請直接輸出 Markdown 格式，不要包含專案名稱標題）：

### 主要工作項目

1. **功能開發**
   - [列出新增的功能]

2. **問題修復**
   - [列出修復的問題]

3. **優化改進**
   - [列出優化和改進項目]

### 技術亮點

- [列出值得注意的技術實作或創新點]

### 明日計畫

- [根據今日進度，建議明日可能的工作重點]

注意：
- 如果某個類別（功能開發/問題修復/優化改進）沒有相關內容，請省略該類別
- 請根據實際 commit 內容分析，不要編造不存在的工作項目
- 技術亮點應該簡潔有力，突出重點
- 明日計畫應該合理且具體`,
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
			MaxCompletionTokens: 1000,
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
	promptBuilder.WriteString("以下是多個專案的今日 commit 訊息，請為每個專案生成結構化的工作總結：\n\n")

	for _, record := range records {
		promptBuilder.WriteString(fmt.Sprintf("## 專案「%s」\n", record.ProjectName))
		for i, msg := range record.CommitMessage {
			promptBuilder.WriteString(fmt.Sprintf("%d. %s\n", i+1, msg))
		}
		promptBuilder.WriteString("\n")
	}

	promptBuilder.WriteString(`請以以下格式用繁體中文回覆，每個專案包含完整的結構化總結：

【專案名稱】
### 主要工作項目

1. **功能開發**
   - [列出新增的功能]

2. **問題修復**
   - [列出修復的問題]

3. **優化改進**
   - [列出優化和改進項目]

### 技術亮點

- [列出值得注意的技術實作或創新點]

### 明日計畫

- [根據今日進度，建議明日可能的工作重點]

注意：
- 每個專案都要以「【專案名稱】」開頭（使用全形括號）
- 如果某個類別（功能開發/問題修復/優化改進）沒有相關內容，請省略該類別
- 請根據實際 commit 內容分析，不要編造不存在的工作項目
- 技術亮點應該簡潔有力，突出重點
- 明日計畫應該合理且具體
`)

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
			MaxCompletionTokens: 2000,
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
