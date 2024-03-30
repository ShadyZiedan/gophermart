package integration

import (
	"context"
	"slices"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"

	"github.com/ShadyZiedan/gophermart/internal/logger"
)

type AccrualService struct {
	client *resty.Client
	orderService
}
type orderService interface {
	SetOrderAccrualResult(ctx context.Context, orderNumber string, status string, accrual float64) error
}

type OrderAccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewAccrualService(address string, orderService orderService) *AccrualService {
	client := resty.New().SetBaseURL(address)
	return &AccrualService{client: client, orderService: orderService}
}

func (a *AccrualService) HandleOrder(ctx context.Context, orderNumber string) error {
	return a.observeOrderAccrualResult(ctx, orderNumber)
}

func (a *AccrualService) getOrderAccrualResult(ctx context.Context, orderNumber string) (*OrderAccrualResponse, error) {
	var response OrderAccrualResponse
	_, err := a.client.R().SetContext(ctx).SetResult(&response).Get(`/api/orders/` + orderNumber)
	if err != nil {
		logger.Log.Error("error getting order accrual", zap.Error(err))
		return nil, err
	}
	return &response, nil
}

func (a *AccrualService) observeOrderAccrualResult(ctx context.Context, orderNumber string) error {
	<-time.After(time.Second)
	var status string
	for {
		accrual, err := a.getOrderAccrualResult(ctx, orderNumber)
		if err != nil {
			logger.Log.Error("error getting order accrual", zap.Error(err))
			return err
		}
		status = accrual.Status
		err = a.orderService.SetOrderAccrualResult(ctx, orderNumber, status, accrual.Accrual)
		if err != nil {
			return err
		}
		if slices.Contains([]string{"PROCESSED", "INVALID"}, status) {
			return nil
		}
	}
}
