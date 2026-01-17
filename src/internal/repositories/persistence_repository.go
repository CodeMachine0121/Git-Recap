package repositories

import "commit-record/src/internal/domains"

type IPersistenceRepository interface {
	Save(record domains.DailyWorkConclusionRecord) error
}

type PersistenceRepository struct{}

func (*PersistenceRepository) Save(domains.DailyWorkConclusionRecord) error {
	return nil
}
