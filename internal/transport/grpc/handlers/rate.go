package handlers

import (
	"GetCurrency/internal/service"
	"GetCurrency/pkg/logger"
	rate_pb "GetCurrency/proto/rate.pb"
	"context"
)

type RateHandler struct {
	service *service.RateService
	logger  logger.Logger
	rate_pb.UnimplementedRateServiceServer
}

func NewRateHandler(service *service.RateService, logger logger.Logger) *RateHandler {
	return &RateHandler{
		service: service,
		logger:  logger,
	}
}

func (h *RateHandler) GetRates(ctx context.Context, req *rate_pb.GetRatesRequest) (*rate_pb.RateResponse, error) {

	rate, err := h.service.GetActualRate(ctx, req.Currency)
	if err != nil {
		h.logger.Errorf("failed to get latest rate:%w", err)
		return nil, err
	}
	h.logger.Debug("raw rate data",
		"ask", rate.Ask,
		"bid", rate.Bid,
		"currency", rate.Currency,
	)

	ask := rate.Ask.InexactFloat64()
	bid := rate.Bid.InexactFloat64()

	response := &rate_pb.RateResponse{
		Id:        rate.ID.String(),
		Currency:  rate.Currency,
		Ask:       ask,
		Bid:       bid,
		Timestamp: rate.Timestamp,
	}
	return response, nil
}
