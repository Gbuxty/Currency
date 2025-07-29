package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"GetCurrency/internal/domain"
	"GetCurrency/internal/service"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRateRepository struct {
	mock.Mock
}

func (m *MockRateRepository) SaveDataRate(ctx context.Context, data *domain.Rate) (*domain.Rate, error) {
	args := m.Called(ctx, data)
	return args.Get(0).(*domain.Rate), args.Error(1)
}


type MockGrinexApi struct {
	mock.Mock
}

func (m *MockGrinexApi) GetCurrencyGrinex(ctx context.Context, currency string) (*domain.Rate, error) {
	args := m.Called(ctx, currency)
	return args.Get(0).(*domain.Rate), args.Error(1)
}

func TestGetActualRate_Success(t *testing.T) {
	ctx := context.Background()
	currency := "usdtrub"

	mockGrinex := new(MockGrinexApi)
	mockRepo := new(MockRateRepository)

	expectedRate := &domain.Rate{
		ID:        uuid.New(),
		Currency:  currency,
		Ask:       decimal.NewFromFloat(101.5),
		Bid:       decimal.NewFromFloat(100.5),
		Timestamp: time.Now().Unix(),
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	mockGrinex.On("GetCurrencyGrinex", ctx, currency).Return(expectedRate, nil)
	mockRepo.On("SaveDataRate", ctx, expectedRate).Return(expectedRate, nil)

	svc := service.NewRateService(mockRepo, mockGrinex)

	rate, err := svc.GetActualRate(ctx, currency)
	assert.NoError(t, err)
	assert.Equal(t, expectedRate, rate)

	mockGrinex.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetActualRate_GrinexError(t *testing.T) {
	ctx := context.Background()
	currency := "usdtrub"

	mockGrinex := new(MockGrinexApi)
	mockRepo := new(MockRateRepository)

	mockGrinex.On("GetCurrencyGrinex", ctx, currency).Return((*domain.Rate)(nil), errors.New("api error"))

	svc := service.NewRateService(mockRepo, mockGrinex)

	rate, err := svc.GetActualRate(ctx, currency)
	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "failed get rate from api")

	mockGrinex.AssertExpectations(t)
}

func TestGetActualRate_SaveDataRateError(t *testing.T) {
	ctx := context.Background()
	currency := "usd"

	mockGrinex := new(MockGrinexApi)
	mockRepo := new(MockRateRepository)

	expectedRate := &domain.Rate{
		ID:        uuid.New(),
		Currency:  currency,
		Ask:       decimal.NewFromFloat(101.5),
		Bid:       decimal.NewFromFloat(100.5),
		Timestamp: time.Now().Unix(),
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	mockGrinex.On("GetCurrencyGrinex", ctx, currency).Return(expectedRate, nil)
	mockRepo.On("SaveDataRate", ctx, expectedRate).Return((*domain.Rate)(nil), errors.New("db error"))

	svc := service.NewRateService(mockRepo, mockGrinex)

	rate, err := svc.GetActualRate(ctx, currency)
	assert.Error(t, err)
	assert.Nil(t, rate)
	assert.Contains(t, err.Error(), "failed save data rate")

	mockGrinex.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}