package postgres

import (
	"context"
	"fmt"

	"github.com/rocky2015aaa/ethdefender/internal/repository/models"
)

func (d *Database) UpsertSlitherReport(ctx context.Context, report *models.ContractReport) error {
	if err := d.Gorm.
		WithContext(ctx).
		Create(&report).Error; err != nil {
		return fmt.Errorf("failed to insert a report: %w", err)
	}

	return nil
}

func (d *Database) GetSlitherReport(ctx context.Context) ([]*models.ContractReport, error) {
	var reports []*models.ContractReport
	// Fetch byte data using GORM
	if err := d.Gorm.WithContext(ctx).
		Find(&reports).
		Error; err != nil {
		return nil, fmt.Errorf("failed to get a report: %w", err)
	}
	return reports, nil
}
