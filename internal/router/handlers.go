package router

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/condratf/go-musthave-diploma-tpl/internal/utils"
)

func (appRouter *AppRouter) uploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := appRouter.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	order := strings.TrimSpace(string(bodyBytes))

	if !regexp.MustCompile(`^\d+$`).MatchString(order) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	if !utils.IsValidLuhn(order) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err = appRouter.ordersRepo.UploadOrder(login, order)
	if err != nil {
		if errors.Is(err, repository.ErrOrderAlreadyUploadedBySameUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, repository.ErrOrderAlreadyUploadedByAnotherUser) {
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
			return
		}
		http.Error(w, "Error uploading order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
func (appRouter *AppRouter) getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := appRouter.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := appRouter.ordersRepo.GetOrders(login)
	if err != nil {
		http.Error(w, "Error getting orders", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
func (appRouter *AppRouter) getBalanceHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := appRouter.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := appRouter.ordersRepo.GetBalance(login)
	if err != nil {
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
func (appRouter *AppRouter) withdrawHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := appRouter.checkSession(r)
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

	if !utils.IsValidLuhn(req.Order) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err := appRouter.ordersRepo.Withdraw(login, req.Order, req.Sum)
	switch {
	case errors.Is(err, repository.ErrInsufficientFunds):
		http.Error(w, "Not enough funds", http.StatusPaymentRequired)
	case err != nil:
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
func (appRouter *AppRouter) getWithdrawalsHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := appRouter.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := appRouter.ordersRepo.GetWithdrawals(login)
	if err != nil {
		http.Error(w, "Error getting withdrawals", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(withdrawals)
}
