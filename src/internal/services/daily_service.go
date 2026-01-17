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

func (s *DailyService) DoDailyWorkConclusion(projectPath string) {

	commitRecord := s.commitService.GetDailyCommitMessages(projectPath)
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

func NewDailyService(commitService *GitCommitService, conclusionService *ConclusionService, persistenceRepo repositories.IPersistenceRepository) *DailyService {
	return &DailyService{commitService, conclusionService, persistenceRepo}
}
