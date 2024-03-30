package handlers

import (
	"context"
	"net/http"

	"github.com/ShadyZiedan/gophermart/internal/models"
)

type Handler struct {
	authService
	orderService
	balanceService
	accrualService
}

type authService interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (string, error)
	NewJWTVerifyMiddleware() func(http.Handler) http.Handler
}

type orderService interface {
	CreateOrder(ctx context.Context, userId int, orderNumber int) (*models.Order, error)
	GetOrders(ctx context.Context, userId int) ([]*models.Order, error)
}

type balanceService interface {
	GetUserBalance(ctx context.Context, userId int) (float64, error)
	GetUserWithdrawalBalance(ctx context.Context, userId int) (float64, error)
	WithdrawOrder(ctx context.Context, userId int, orderNumber int, sum float64) error
	GetWithdrawals(ctx context.Context, userId int) ([]*models.Withdrawal, error)
}

type accrualService interface {
	HandleOrder(ctx context.Context, orderNumber string) error
}

func NewHandler(
	authService authService,
	orderService orderService,
	balanceService balanceService,
	accrualService accrualService,
) *Handler {
	return &Handler{
		authService:    authService,
		orderService:   orderService,
		balanceService: balanceService,
		accrualService: accrualService,
	}
}
