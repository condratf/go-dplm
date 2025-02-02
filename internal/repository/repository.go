package repository

import (
	"database/sql"
	"errors"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type Repository interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]string, error)
	GetBalance(login string) (int, error)
	Withdraw(login string, amount int) error
	GetWithdrawals(login string) ([]int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) UploadOrder(login, order string) error {
	_, err := r.db.Exec("INSERT INTO orders (login, order) VALUES ($1, $2)", login, order)
	return err
}

func (r *repository) GetOrders(login string) ([]string, error) {
	rows, err := r.db.Query("SELECT order FROM orders WHERE login = $1", login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []string
	for rows.Next() {
		var order string
		if err := rows.Scan(&order); err != nil {
			return nil, err
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
