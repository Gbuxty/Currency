package repository

import (
	"GetCurrency/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveDataRate(ctx context.Context, data *domain.Rate) (*domain.Rate, error) {
	query := `
        INSERT INTO rates (currency, ask, bid, timestamp) 
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (currency) DO UPDATE 
        SET ask = EXCLUDED.ask, 
            bid = EXCLUDED.bid, 
            timestamp = EXCLUDED.timestamp,
			updated_at = NOW()
		RETURNING id,currency,ask,bid,timestamp
    `
	var freshData domain.Rate
	err := s.db.QueryRow(ctx,
		query,
		data.Currency,
		data.Ask,
		data.Bid,
		data.Timestamp,
		).Scan(
		&freshData.ID,
		&freshData.Currency,
		&freshData.Ask,
		&freshData.Bid,
		&freshData.Timestamp,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to save rate data: %w", err)
	}

	return &freshData, nil
}
