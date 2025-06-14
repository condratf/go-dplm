package withdrawrouter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/errors_custom"
)

func NewWithdrawRouter(
	checkSession func(r *http.Request) (string, bool),
	withdrawService WithdrawService,
) WithdrawRouter {
	return &withdrawRouter{
		checkSession:    checkSession,
		withdrawService: withdrawService,
	}
}

func (router *withdrawRouter) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := router.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := router.withdrawService.Withdraw(login, req.Order, req.Sum)
	switch {
	case errors.Is(err, errors_custom.ErrInvalidOrderNumber):
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
	case errors.Is(err, errors_custom.ErrInsufficientFunds):
		http.Error(w, "Not enough funds", http.StatusPaymentRequired)
	case err != nil:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func (router *withdrawRouter) GetWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := router.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := router.withdrawService.GetWithdrawals(login)
	if err != nil {
		if errors.Is(err, errors_custom.ErrNoContent) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "Error getting withdrawals", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawals)
}
