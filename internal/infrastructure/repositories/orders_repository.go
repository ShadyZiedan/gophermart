package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/ShadyZiedan/gophermart/internal/models"
)

type OrdersRepository struct {
	conn pgConn
}

func NewOrdersRepository(conn pgConn) *OrdersRepository {
	return &OrdersRepository{conn}
}

func (r *OrdersRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	row := r.conn.QueryRow(ctx, "INSERT INTO orders(user_id, number, status) VALUES($1, $2, $3) returning id, uploaded_at", order.UserId, order.Number, "NEW")
	err := row.Scan(&order.Id, &order.UploadedAt)
	return err
}

func (r *OrdersRepository) FindOrderByNumber(ctx context.Context, orderNumber int) (*models.Order, error) {
	var order models.Order
	row := r.conn.QueryRow(ctx, "SELECT id, user_id, number, status, uploaded_at, processed_at FROM orders WHERE number = $1", orderNumber)
	var uploadedAt pgtype.Timestamp
	var processedAt pgtype.Timestamp
	err := row.Scan(&order.Id, &order.UserId, &order.Number, &order.Status, &uploadedAt, &processedAt)
	if err != nil {
		return nil, err
	}
	if uploadedAt.Valid {
		order.UploadedAt = uploadedAt.Time
	}
	if processedAt.Valid {
		order.ProcessedAt = processedAt.Time
	}
	return &order, nil
}

func (r *OrdersRepository) FindOrdersByUserId(ctx context.Context, userId int) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)
	rows, err := r.conn.Query(ctx, "SELECT id, user_id, number, status, accrual, uploaded_at, processed_at FROM orders WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		var uploadedAt pgtype.Timestamp
		var processedAt pgtype.Timestamp
		var accrual pgtype.Float8
		err := rows.Scan(&order.Id, &order.UserId, &order.Number, &order.Status, &accrual, &uploadedAt, &processedAt)
		if err != nil {
			return nil, err
		}
		if accrual.Valid {
			order.Accrual = accrual.Float64
		}
		if uploadedAt.Valid {
			order.UploadedAt = uploadedAt.Time
		}
		if processedAt.Valid {
			order.ProcessedAt = processedAt.Time
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (r *OrdersRepository) UpdateOrderAccrual(ctx context.Context, number string, status string, accrual float64, processedAt time.Time) error {
	if !processedAt.IsZero() {
		sql := `update orders set status = $1, accrual = $2, processed_at = $3 where number = $4`
		_, err := r.conn.Exec(ctx, sql, status, accrual, processedAt, number)
		return err
	}
	sql := `update orders set status = $1, accrual = $2 where number = $3`
	_, err := r.conn.Exec(ctx, sql, status, accrual, number)
	return err

}
