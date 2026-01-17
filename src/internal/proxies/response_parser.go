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
		line = strings.TrimSpace(line)

		if projectName, isHeader := parser.extractProjectName(line); isHeader {
			parser.saveCurrentProject(currentProject, &currentContent, result)
			currentProject = projectName
		} else if line != "" && currentProject != "" {
			currentContent.WriteString(line)
			currentContent.WriteString(" ")
		}
	}

	parser.saveCurrentProject(currentProject, &currentContent, result)
	parser.fillMissingProjects(result, records, content)

	return result
}

func (parser *MultiProjectResponseParser) extractProjectName(line string) (string, bool) {
	if strings.HasPrefix(line, "【") && strings.Contains(line, "】") {
		name := strings.TrimPrefix(line, "【")
		name = strings.TrimSuffix(name, "】")
		return name, true
	}

	if strings.HasPrefix(line, "##") {
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
