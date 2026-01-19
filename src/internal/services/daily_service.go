package services

import (
	"commit-record/src/internal/domains"
	"commit-record/src/internal/repositories"
)

type DailyService struct {
	commitService     *GitCommitService
	conclusionService *ConclusionService
	persistenceRepo   repositories.IPersistenceRepository
}

func (s *DailyService) DoDailyWorkConclusion(author, projectPath string) {

	commitRecord := s.commitService.GetDailyCommitMessages(author, projectPath)
	if commitRecord.CommitMessage == nil || len(commitRecord.CommitMessage) == 0 {
		println("No commit messages found for the project")
		return
	}
	conclusion := s.conclusionService.GetConclusion(commitRecord)
	err := s.persistenceRepo.Save(domains.DailyWorkConclusionRecord{
		ProjectName: commitRecord.ProjectName,
		Conclusion:  conclusion,
	})

	if err != nil {
		panic(err)
	}
}

func (s *DailyService) DoDailyWorkConclusionBatch(projectPaths []string) {
	var records []domains.CommitRecord

	for _, projectPath := range projectPaths {
		commitRecord := s.commitService.GetDailyCommitMessages(projectPath, "")
		if commitRecord.CommitMessage != nil && len(commitRecord.CommitMessage) > 0 {
			records = append(records, commitRecord)
		}
	}

	if len(records) == 0 {
		println("No commit messages found for any project")
		return
	}

	conclusions := s.conclusionService.GetBatchConclusion(records)

	for _, record := range records {
		conclusion, exists := conclusions[record.ProjectName]
		if !exists {
			println("Warning: No conclusion found for project:", record.ProjectName)
			continue
		}

		err := s.persistenceRepo.Save(domains.DailyWorkConclusionRecord{
			ProjectName: record.ProjectName,
			Conclusion:  conclusion,
		})

		if err != nil {
			panic(err)
		}
	}
}

func NewDailyService(commitService *GitCommitService, conclusionService *ConclusionService, persistenceRepo repositories.IPersistenceRepository) *DailyService {
	return &DailyService{commitService, conclusionService, persistenceRepo}
}
