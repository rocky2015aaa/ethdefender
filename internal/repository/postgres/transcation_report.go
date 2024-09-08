package postgres

import (
	"context"
	"fmt"

	"github.com/rocky2015aaa/ethdefender/internal/repository/models"
)

func (d *Database) CreateTransactionReport(ctx context.Context, report *models.TransactionReport) error {
	if err := d.Gorm.
		WithContext(ctx).
		Create(&report).Error; err != nil {
		return fmt.Errorf("failed to insert a report: %w", err)
	}

	return nil
}

func (d *Database) GetTransactionReport(ctx context.Context) ([]*models.TransactionReport, error) {
	var reports []*models.TransactionReport

	if err := d.Gorm.
		WithContext(ctx).
		Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	return reports, nil
}
