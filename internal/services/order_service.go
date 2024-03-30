package services

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/ShadyZiedan/gophermart/internal/luhn"
	"github.com/ShadyZiedan/gophermart/internal/models"
)

type OrderService struct {
	orderRepository
}

func NewOrderService(orderRepository orderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

type orderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	FindOrderByNumber(ctx context.Context, orderNumber int) (*models.Order, error)
	FindOrdersByUserId(ctx context.Context, userID int) ([]*models.Order, error)
	UpdateOrderAccrual(ctx context.Context, number string, status string, accrual float64, processedAt time.Time) error
}

var (
	ErrOrderAlreadyExists      = errors.New("order already exists")
	ErrOrderBelongsToOtherUser = errors.New("order belongs to other user")
	ErrInvalidOrderNumber      = errors.New("invalid order number")
)

func (o *OrderService) CreateOrder(ctx context.Context, userID int, orderNumber int) (*models.Order, error) {
	order, err := o.orderRepository.FindOrderByNumber(ctx, orderNumber)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if order != nil {
		if order.UserID != userID {
			return nil, ErrOrderBelongsToOtherUser
		} else {
			return nil, ErrOrderAlreadyExists
		}
	}

	if !luhn.Valid(orderNumber) {
		return nil, ErrInvalidOrderNumber
	}

	order = &models.Order{UserID: userID, Number: orderNumber}
	err = o.orderRepository.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *OrderService) GetOrders(ctx context.Context, userID int) ([]*models.Order, error) {
	return o.orderRepository.FindOrdersByUserId(ctx, userID)
}

func (o *OrderService) SetOrderAccrualResult(ctx context.Context, orderNumber string, status string, accrual float64) error {
	var processedAt time.Time
	if slices.Contains([]string{"PROCESSED", "INVALID"}, status) {
		processedAt = time.Now()
	}
	return o.orderRepository.UpdateOrderAccrual(ctx, orderNumber, status, accrual, processedAt)
}
