package withdrawservice

import (
	"github.com/condratf/go-musthave-diploma-tpl/internal/errors_custom"
	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
	"github.com/condratf/go-musthave-diploma-tpl/internal/utils"
)

func NewWithdrawService(
	withdrawRepo withdrawRepository,
) WithdrawService {
	return &withdrawService{
		withdrawRepo: withdrawRepo,
	}
}

func (w *withdrawService) Withdraw(login string, order string, amount float64) error {
	if !utils.IsValidLuhn(order) {
		return errors_custom.ErrInvalidOrderNumber
	}
	err := w.withdrawRepo.Withdraw(login, order, amount)
	return err
}

func (w *withdrawService) GetWithdrawals(login string) ([]models.WithdrawalRes, error) {
	withdrawals, err := w.withdrawRepo.GetWithdrawals(login)
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, errors_custom.ErrNoContent
	}

	return withdrawals, nil
}
