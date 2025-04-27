package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrUserExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(login, password string) error
	GetUserPassword(login string) (string, error)
	GetUserBalance(login string) (int, error)
	UpdateUserBalance(login string, amount int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(login, password string) error {
	_, err := r.db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrUserExists
		}
		return err
	}
	return nil
}

func (r *userRepository) GetUserPassword(login string) (string, error) {
	var password string
	err := r.db.QueryRow("SELECT password FROM users WHERE login = $1", login).Scan(&password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}
	return password, nil
}

func (r *userRepository) GetUserBalance(login string) (int, error) {
	var balance int
	err := r.db.QueryRow("SELECT balance FROM users WHERE login = $1", login).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (r *userRepository) UpdateUserBalance(login string, amount int) error {
	_, err := r.db.Exec("UPDATE users SET balance = balance - $1 WHERE login = $2", amount, login)
	return err
}
