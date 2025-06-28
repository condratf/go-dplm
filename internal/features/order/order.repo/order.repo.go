package orderrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/condratf/go-musthave-diploma-tpl/internal/custerrors"
	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) UploadOrder(login, order string) error {
	var existingLogin string
	err := r.db.QueryRow("SELECT user_login FROM orders WHERE order_number = $1", order).Scan(&existingLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.db.Exec("INSERT INTO orders (user_login, order_number) VALUES ($1, $2)", login, order)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if existingLogin == login {
		return custerrors.ErrOrderAlreadyUploadedBySameUser
	}
	return custerrors.ErrOrderAlreadyUploadedByAnotherUser
}

func (r *orderRepository) GetOrders(login string) ([]models.Order, error) {
	rows, err := r.db.Query(`
			SELECT order_number, status, loyalty_points, created_at 
			FROM orders 
			WHERE user_login = $1 
			ORDER BY created_at DESC; 
	`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		var createdAt time.Time
		var accrual sql.NullInt64

		if err := rows.Scan(&order.Order, &order.Status, &accrual, &createdAt); err != nil {
			return nil, err
		}

		order.CreatedAt = createdAt.Format(time.RFC3339)
		if accrual.Valid {
			loyaltyPoints := int(accrual.Int64)
			order.LoyaltyPoints = loyaltyPoints
		}

		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
			SELECT order_number 
			FROM orders
			WHERE status IN ('pending', 'PENDING')
			ORDER BY created_at LIMIT 100;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.Order); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error {
	_, err := r.db.ExecContext(ctx, `
			UPDATE orders 
			SET status = $1, loyalty_points = $2 
			WHERE order_number = $3
	`, status, accrual, orderNumber)
	return err
}

func (r *orderRepository) GetBalance(login string) (*models.BalanceResponse, error) {
	balance := &models.BalanceResponse{}

	query := `
			SELECT ub.available_points, COALESCE(SUM(w.amount), 0)
			FROM users u
			JOIN user_balance ub ON u.id = ub.user_id
			LEFT JOIN withdrawals w ON u.id = w.user_id
			WHERE u.login = $1
			GROUP BY ub.available_points
	`
	err := r.db.QueryRow(query, login).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
