package services

import (
	domains "commit-record/src/internal/domains"
	"commit-record/src/internal/git"
)

type DailyService struct {
	gitHandler git.IGitHandler
}

func (s *DailyService) GetDailyCommitMessages(projectPath string) domains.CommitRecord {

	commitRecord := s.gitHandler.GetCommitMessages(projectPath)

	return domains.CommitRecord{ProjectName: projectPath, CommitMessage: commitRecord}
}

func NewDailyService(gitHandler git.IGitHandler) *DailyService {
	return &DailyService{gitHandler}
}
