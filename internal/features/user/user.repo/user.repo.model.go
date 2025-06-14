package userrepo

import "database/sql"

type UserRepository interface {
	CreateUser(login, password, email string) error
	GetUserPassword(login string) (string, error)
	GetUserBalance(login string) (int, error)
	UpdateUserBalance(login string, amount int) error
}

type userRepository struct {
	db *sql.DB
}
