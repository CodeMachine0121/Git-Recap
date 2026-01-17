package services

import (
	"commit-record/src/internal/domains"
	"commit-record/src/internal/git"
)

type GitCommitService struct {
	gitHandler git.IGitHandler
}

func (s *GitCommitService) GetDailyCommitMessages(projectPath string) domains.CommitRecord {

	commitRecord := s.gitHandler.GetCommitMessages(projectPath)

	return domains.CommitRecord{ProjectName: projectPath, CommitMessage: commitRecord}
}

func NewGitCommitService(gitHandler git.IGitHandler) *GitCommitService {
	return &GitCommitService{gitHandler}
}
