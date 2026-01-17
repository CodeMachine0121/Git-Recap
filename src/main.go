package main

import (
	"flag"
	"fmt"
	"os"

	"commit-record/src/internal/git"
	"commit-record/src/internal/proxies"
	"commit-record/src/internal/repositories"
	"commit-record/src/internal/services"
)

func main() {
	flag.Parse()
	projectPaths := flag.Args()

	if len(projectPaths) == 0 {
		fmt.Println("Usage: commit-record <project-path-1> [<project-path-2> ...]")
		os.Exit(1)
	}

	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable is required")
		os.Exit(1)
	}

	gitHandler := &git.GitHandler{}
	persistenceRepo := &repositories.PersistenceRepository{}
	openAiProxy := proxies.NewOpenAiProxy()

	gitService := services.NewGitCommitService(gitHandler)
	conclusionService := services.NewConclusionService(openAiProxy)
	dailyService := services.NewDailyService(gitService, conclusionService, persistenceRepo)

	// 驗證所有專案路徑
	for _, projectPath := range projectPaths {
		if _, err := os.Stat(projectPath); os.IsNotExist(err) {
			fmt.Printf("Error: Path does not exist: %s\n", projectPath)
			os.Exit(1)
		}
	}

	dailyService.DoDailyWorkConclusionBatch(projectPaths)
}
