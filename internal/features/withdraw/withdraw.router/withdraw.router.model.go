package withdrawrouter

import (
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type WithdrawRouter interface {
	WithdrawHandler(w http.ResponseWriter, r *http.Request)
	GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request)
}

type WithdrawService interface {
	GetWithdrawals(login string) ([]models.WithdrawalRes, error)
	Withdraw(login string, order string, amount float64) error
}

type withdrawRouter struct {
	checkSession    func(r *http.Request) (string, bool)
	withdrawService WithdrawService
}
