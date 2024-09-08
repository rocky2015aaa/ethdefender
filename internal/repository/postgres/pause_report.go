package postgres

import (
	"context"
	"fmt"

	"github.com/rocky2015aaa/ethdefender/internal/repository/models"
)

func (d *Database) CreatePauseReport(ctx context.Context, report *models.PauseReport) error {
	if err := d.Gorm.
		WithContext(ctx).
		Create(&report).Error; err != nil {
		return fmt.Errorf("failed to insert a report: %w", err)
	}

	return nil
}

func (d *Database) GetPauseReport(ctx context.Context) ([]*models.PauseReport, error) {
	var reports []*models.PauseReport

	if err := d.Gorm.
		WithContext(ctx).
		Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	return reports, nil
}
