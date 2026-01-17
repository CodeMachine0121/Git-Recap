package proxies

import "commit-record/src/internal/domains"

type IOpenAiProxy interface {
	GetConclusion(record domains.CommitRecord) string
}
