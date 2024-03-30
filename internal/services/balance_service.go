package services

import (
	"context"
	"errors"

	"github.com/ShadyZiedan/gophermart/internal/infrastructure/repositories"
	"github.com/ShadyZiedan/gophermart/internal/models"
)

type BalanceService struct {
	balanceRepository
}

func NewBalanceService(balanceRepository balanceRepository) *BalanceService {
	return &BalanceService{balanceRepository: balanceRepository}
}

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type balanceRepository interface {
	GetUserBalance(ctx context.Context, userId int) (float64, error)
	GetUserWithdrawalBalance(ctx context.Context, userId int) (float64, error)
	CreateWithdrawal(ctx context.Context, userId int, orderNumber int, sum float64) (*models.Withdrawal, error)
	GetWithdrawals(ctx context.Context, userId int) ([]*models.Withdrawal, error)
}

func (b *BalanceService) WithdrawOrder(ctx context.Context, userId int, orderNumber int, sum float64) error {
	_, err := b.balanceRepository.CreateWithdrawal(ctx, userId, orderNumber, sum)
	if err != nil {
		if errors.Is(err, repositories.ErrInsufficientBalance) {
			return ErrInsufficientBalance
		}
		return err
	}
	return nil
}

func (b *BalanceService) GetUserBalance(ctx context.Context, userId int) (float64, error) {
	balance, err := b.balanceRepository.GetUserBalance(ctx, userId)
	return balance, err
}

func (b *BalanceService) GetUserWithdrawalBalance(ctx context.Context, userId int) (float64, error) {
	balance, err := b.balanceRepository.GetUserWithdrawalBalance(ctx, userId)
	return balance, err
}
