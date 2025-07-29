package service

import (
	"GetCurrency/internal/domain"
	"context"
	"fmt"
)

type RateService struct {
	repo   RateRepository
	grinex GrinexApi
}

type RateRepository interface {
	SaveDataRate(ctx context.Context, data *domain.Rate) (*domain.Rate, error)
}

type GrinexApi interface {
	GetCurrencyGrinex(ctx context.Context, currency string) (*domain.Rate, error)
}

func NewRateService(repo RateRepository, grinex GrinexApi) *RateService {
	return &RateService{
		repo:   repo,
		grinex: grinex,
	}
}

func (s *RateService) GetActualRate(ctx context.Context, currency string) (*domain.Rate, error) {
	rate, err := s.grinex.GetCurrencyGrinex(ctx, currency)
	if err != nil {
		return nil, fmt.Errorf("failed get rate from api:%w", err)
	}

	newData, err := s.repo.SaveDataRate(ctx, rate)
	if err != nil {
		return nil, fmt.Errorf("failed save data rate:%w", err)
	}

	return newData, nil
}
