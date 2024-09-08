package models

import (
	"time"
)

type ContractReport struct {
	ID           uint      `gorm:"primaryKey"`
	ContractName string    `json:"contract_name"`
	Report       []byte    `gorm:"type:bytea"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TransactionReport struct {
	ID            uint      `gorm:"primaryKey"`
	TransactionID string    `json:"transaction_id"`
	ExecutionTime string    `json:"execution_time"`
	GasUsed       uint64    `json:"gas_used"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PauseReport struct {
	ID          uint      `gorm:"primaryKey"`
	EventType   string    `json:"event_type"`
	PauseStatus bool      `json:"pause_status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
