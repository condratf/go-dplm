package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

var (
	ErrOrderAlreadyUploadedBySameUser    = errors.New("order already uploaded by same user")
	ErrOrderAlreadyUploadedByAnotherUser = errors.New("order already uploaded by another user")
	ErrInsufficientFunds                 = errors.New("insufficient funds")
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

type OrdersRepository interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetBalance(login string) (BalanceResponse, error)
	Withdraw(login string, order string, amount float64) error
	GetWithdrawals(login string) ([]Withdrawal, error)
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error
}
type repository struct {
	db *sql.DB
}

func NewOrdersRepository(db *sql.DB) OrdersRepository {
	return &repository{db: db}
}

func (r *repository) UploadOrder(login, order string) error {
	var existingLogin string
	err := r.db.QueryRow("SELECT user_login FROM orders WHERE order = $1", order).Scan(&existingLogin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.db.Exec("INSERT INTO orders (login, order) VALUES ($1, $2)", login, order)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if existingLogin == login {
		return ErrOrderAlreadyUploadedBySameUser
	}
	return ErrOrderAlreadyUploadedByAnotherUser
}

func (r *repository) GetOrders(login string) ([]models.Order, error) {
	rows, err := r.db.Query(`
			SELECT order, status, loyalty_points, created_at 
			FROM orders 
			WHERE user_login = $1 
			ORDER BY created_at DESC
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
	return orders, nil
}

func (r *repository) GetBalance(login string) (BalanceResponse, error) {
	var balance BalanceResponse

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
		return BalanceResponse{}, err
	}

	return balance, nil
}

func (r *repository) Withdraw(login, order string, amount float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	var available float64

	err = tx.QueryRow(`
		SELECT u.id, ub.available_points 
		FROM users u
		JOIN user_balance ub ON u.id = ub.user_id
		WHERE u.login = $1
	`, login).Scan(&userID, &available)
	if err != nil {
		return err
	}

	if available < amount {
		return ErrInsufficientFunds
	}

	_, err = tx.Exec(`
		UPDATE user_balance 
		SET available_points = available_points - $1 
		WHERE user_id = $2
	`, amount, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO withdrawals (user_id, amount, order_number, created_at) 
		VALUES ($1, $2, $3, NOW())
	`, userID, amount, order)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetWithdrawals(login string) ([]Withdrawal, error) {
	query := `
		SELECT w.order_number, w.amount, w.created_at
		FROM users u
		JOIN withdrawals w ON u.id = w.user_id
		WHERE u.login = $1
		ORDER BY w.created_at DESC
	`
	rows, err := r.db.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []Withdrawal
	for rows.Next() {
		var w Withdrawal
		var createdAt time.Time
		if err := rows.Scan(&w.Order, &w.Sum, &createdAt); err != nil {
			return nil, err
		}
		w.ProcessedAt = createdAt.Format(time.RFC3339)
		withdrawals = append(withdrawals, w)
	}
	return withdrawals, nil
}

func (r *repository) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
			SELECT order FROM orders 
			WHERE status IN ('REGISTERED', 'PROCESSING')
			ORDER BY created_at LIMIT 100
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
	return orders, nil
}

func (r *repository) UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error {
	_, err := r.db.ExecContext(ctx, `
			UPDATE orders 
			SET status = $1, loyalty_points = $2 
			WHERE order = $3
	`, status, accrual, orderNumber)
	return err
}
