package repository

import (
	"context"

	"github.com/rocky2015aaa/ethdefender/internal/repository/models"
)

type Storage interface {
	CreatePauseReport(ctx context.Context, report *models.PauseReport) error
	GetPauseReport(ctx context.Context) ([]*models.PauseReport, error)

	UpsertSlitherReport(ctx context.Context, report *models.ContractReport) error
	GetSlitherReport(ctx context.Context) ([]*models.ContractReport, error)

	CreateTransactionReport(ctx context.Context, report *models.TransactionReport) error
	GetTransactionReport(ctx context.Context) ([]*models.TransactionReport, error)
}
