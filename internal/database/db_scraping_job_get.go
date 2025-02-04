package database

import (
	"context"
	"fmt"

	lead_scraper_servicev1 "github.com/VectorEngineering/vector-protobuf-definitions/api-definitions/pkg/generated/lead_scraper_service/v1"
)

// GetScrapingJob retrieves a scraping job by ID
func (db *Db) GetScrapingJob(ctx context.Context, id uint64) (*lead_scraper_servicev1.ScrapingJob, error) {
	if id == 0 {
		return nil, ErrInvalidInput
	}

	ctx, cancel := context.WithTimeout(ctx, db.GetQueryTimeout())
	defer cancel()

	var jobORM lead_scraper_servicev1.ScrapingJobORM
	result := db.Client.Engine.WithContext(ctx).Where("id = ?", id).First(&jobORM)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, fmt.Errorf("%w: %v", ErrJobDoesNotExist, result.Error)
		}
		return nil, fmt.Errorf("failed to get scraping job: %w", result.Error)
	}

	pbResult, err := jobORM.ToPB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to protobuf: %w", err)
	}

	return &pbResult, nil
}
