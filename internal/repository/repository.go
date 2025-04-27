package repository

import (
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

type OrdersRepository interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetBalance(login string) (int, error)
	Withdraw(login string, amount int) error
	GetWithdrawals(login string) ([]int, error)
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

func (r *repository) GetBalance(login string) (int, error) {
	var balance int
	err := r.db.QueryRow("SELECT balance FROM users WHERE login = $1", login).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *repository) Withdraw(login string, amount int) error {
	var balance int
	err := r.db.QueryRow("SELECT balance FROM users WHERE login = $1", login).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return ErrInsufficientFunds
	}

	_, err = r.db.Exec("UPDATE users SET balance = balance - $1 WHERE login = $2", amount, login)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("INSERT INTO withdrawals (login, amount) VALUES ($1, $2)", login, amount)
	return err
}

func (r *repository) GetWithdrawals(login string) ([]int, error) {
	rows, err := r.db.Query("SELECT amount FROM withdrawals WHERE login = $1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []int
	for rows.Next() {
		var amount int
		if err := rows.Scan(&amount); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, amount)
	}
	return withdrawals, nil
}
