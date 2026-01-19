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
	benchMode := flag.Bool("bench", false, "使用批次處理模式")
	singleMode := flag.Bool("single", false, "使用單一專案處理模式")
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Usage: commit-record [--bench|--single] <author-name> <project-path-1> [<project-path-2> ...]")
		os.Exit(1)
	}

	author := args[0]
	projectPaths := args[1:]

	if *benchMode && *singleMode {
		fmt.Println("Error: Cannot use both --bench and --single flags")
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

	// 根據標誌選擇執行模式
	if *benchMode {
		dailyService.DoDailyWorkConclusionBatch(author, projectPaths)
	} else if *singleMode {
		for _, projectPath := range projectPaths {
			dailyService.DoDailyWorkConclusion(author, projectPath)
		}
	} else {
		// 預設使用批次模式
		dailyService.DoDailyWorkConclusionBatch(author, projectPaths)
	}
}
