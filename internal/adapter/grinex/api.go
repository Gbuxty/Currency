package grinex

import (
	"GetCurrency/internal/config"
	"GetCurrency/internal/domain"
	"GetCurrency/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type Grinex struct {
	client *http.Client
	cfg    *config.GrinexConfig
	logger logger.Logger
}

func NewGrinex(cfg *config.GrinexConfig, logger logger.Logger) *Grinex {
	return &Grinex{
		cfg:    cfg,
		client: &http.Client{},
		logger: logger,
	}
}

type Order struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type RateResp struct {
	Timestamp int64   `json:"timestamp"`
	Asks      []Order `json:"asks"`
	Bids      []Order `json:"bids"`
}

func (g *Grinex) GetCurrencyGrinex(ctx context.Context, currency string) (*domain.Rate, error) {
	url := fmt.Sprintf("%s%s", g.cfg.Url, currency)
	g.logger.Debug("requesting rate from Grinex", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			g.logger.Errorf("failed close body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rate RateResp

	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(rate.Asks) == 0 || len(rate.Bids) == 0 {
		return nil, fmt.Errorf("no data in response")
	}

	ask, err := decimal.NewFromString(rate.Asks[0].Price)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ask price: %w", err)

	}

	bid, err := decimal.NewFromString(rate.Bids[0].Price)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bid price: %w", err)
	}

	g.logger.Debug("received rate data from Grinex", ask, bid)

	return &domain.Rate{
		Currency:  currency,
		Ask:       ask,
		Bid:       bid,
		Timestamp: rate.Timestamp,
	}, nil
}
