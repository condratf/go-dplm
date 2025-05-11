package withdrawservice

import "github.com/condratf/go-musthave-diploma-tpl/internal/models"

type WithdrawService interface {
	Withdraw(login string, order string, amount float64) error
	GetWithdrawals(login string) ([]models.WithdrawalRes, error)
}

type withdrawRepository interface {
	Withdraw(login string, order string, amount float64) error
	GetWithdrawals(login string) ([]models.WithdrawalRes, error)
}

type withdrawService struct {
	withdrawRepo withdrawRepository
}
