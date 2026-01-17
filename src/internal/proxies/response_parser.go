package proxies

import (
	"commit-record/src/internal/domains"
	"strings"
)

type MultiProjectResponseParser struct{}

func NewMultiProjectResponseParser() *MultiProjectResponseParser {
	return &MultiProjectResponseParser{}
}

func (parser *MultiProjectResponseParser) Parse(content string, records []domains.CommitRecord) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(content, "\n")
	currentProject := ""
	var currentContent strings.Builder

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if projectName, isHeader := parser.extractProjectName(trimmedLine); isHeader {
			parser.saveCurrentProject(currentProject, &currentContent, result)
			currentProject = projectName
		} else if currentProject != "" {
			// 保留原始行（包括前導空格以保持縮排）但移除行尾空白
			currentContent.WriteString(strings.TrimRight(line, " \t"))
			currentContent.WriteString("\n")
		}
	}

	parser.saveCurrentProject(currentProject, &currentContent, result)
	parser.fillMissingProjects(result, records, content)

	return result
}

func (parser *MultiProjectResponseParser) extractProjectName(line string) (string, bool) {
	// 匹配 【專案名稱】 格式（這是我們在 prompt 中要求的格式）
	if strings.HasPrefix(line, "【") && strings.Contains(line, "】") {
		name := strings.TrimPrefix(line, "【")
		name = strings.TrimSuffix(name, "】")
		return name, true
	}

	// 只匹配 "## " 開頭（兩個井號加空格），避免匹配 "###"（三個井號）等內文標題
	// 同時要確保不是 "### " 開頭
	if strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "### ") {
		name := strings.TrimPrefix(line, "## ")
		name = strings.TrimPrefix(name, "專案「")
		name = strings.TrimSuffix(name, "」")
		return name, true
	}

	return "", false
}

func (parser *MultiProjectResponseParser) saveCurrentProject(projectName string, content *strings.Builder, result map[string]string) {
	if projectName != "" {
		result[projectName] = strings.TrimSpace(content.String())
		content.Reset()
	}
}

func (parser *MultiProjectResponseParser) fillMissingProjects(result map[string]string, records []domains.CommitRecord, fallbackContent string) {
	for _, record := range records {
		if _, exists := result[record.ProjectName]; !exists {
			result[record.ProjectName] = fallbackContent
		}
	}
}
