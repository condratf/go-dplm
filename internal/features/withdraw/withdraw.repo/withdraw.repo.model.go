package withdrawrepo

import (
	"database/sql"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type WithdrawRepository interface {
	Withdraw(login, order string, amount float64) error
	GetWithdrawals(login string) ([]models.WithdrawalRes, error)
}

type withdrawRepository struct {
	db *sql.DB
}
