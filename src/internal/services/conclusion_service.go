package services

import (
	"commit-record/src/internal/domains"
	"commit-record/src/internal/proxies"
)

type ConclusionService struct {
	proxy proxies.IOpenAiProxy
}

func (s *ConclusionService) GetConclusion(record domains.CommitRecord) string {
	conclusion := s.proxy.GetConclusion(record)
	return conclusion
}

func NewConclusionService(proxy proxies.IOpenAiProxy) *ConclusionService {
	return &ConclusionService{proxy}
}
