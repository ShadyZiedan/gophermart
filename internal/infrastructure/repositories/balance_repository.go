package repositories

import (
	"context"
	"errors"

	"github.com/ShadyZiedan/gophermart/internal/models"
)

type BalanceRepository struct {
	conn pgConn
}

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func (b BalanceRepository) GetUserBalance(ctx context.Context, userId int) (float64, error) {
	sql := `
		-- getting users balance
		select coalesce((SELECT sum(accrual) FROM orders WHERE orders.user_id = $1), 0) - coalesce((select sum(sum) from withdrawals where user_id = $1), 0);
       `
	row := b.conn.QueryRow(ctx, sql, userId)
	var balance float64
	err := row.Scan(&balance)
	return balance, err
}

func (b BalanceRepository) GetUserWithdrawalBalance(ctx context.Context, userId int) (float64, error) {
	sql := `select coalesce(sum(sum), 0) from withdrawals where user_id = $1;`
	row := b.conn.QueryRow(ctx, sql, userId)
	var balance float64
	err := row.Scan(&balance)
	return balance, err
}

func (b BalanceRepository) CreateWithdrawal(ctx context.Context, userId int, orderNumber int, sum float64) (*models.Withdrawal, error) {
	tx, err := b.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var balance float64
	err = tx.QueryRow(ctx, `select sum(accrual) from orders where user_id = $1`, userId).Scan(&balance)
	if err != nil {
		return nil, err
	}
	if balance < sum {
		return nil, ErrInsufficientBalance
	}
	var withdrawal models.Withdrawal
	withdrawal.UserId = userId
	withdrawal.OrderNumber = orderNumber
	withdrawal.Sum = sum

	sql := `insert into withdrawals(user_id, number, sum, processed_at) values($1, $2, $3, current_timestamp) returning id, processed_at;`
	row := tx.QueryRow(ctx, sql, userId, orderNumber, sum)
	err = row.Scan(&withdrawal.Id, &withdrawal.ProcessedAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &withdrawal, nil

}

func (b BalanceRepository) GetWithdrawals(ctx context.Context, userId int) ([]*models.Withdrawal, error) {
	var withdraws []*models.Withdrawal
	sql := `select id, user_id, number, sum, processed_at from withdrawals where user_id = $1;`
	rows, err := b.conn.Query(ctx, sql, userId)
	if err != nil {
		return withdraws, err
	}
	defer rows.Close()
	for rows.Next() {
		var w models.Withdrawal
		err := rows.Scan(&w.Id, &w.UserId, &w.OrderNumber, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, &w)
	}
	return withdraws, nil
}

func NewBalanceRepository(conn pgConn) *BalanceRepository {
	return &BalanceRepository{conn: conn}
}
