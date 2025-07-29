package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Rate struct {
	ID        uuid.UUID
	Currency  string
	Ask       decimal.Decimal
	Bid       decimal.Decimal
	Timestamp int64
	CreatedAt time.Time
	UpdatedAt *time.Time
}
