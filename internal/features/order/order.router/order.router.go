package orderrouter

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/condratf/go-musthave-diploma-tpl/internal/errors_custom"
)

func NewOrderRouter(
	checkSession func(r *http.Request) (string, bool),
	orderService orderService,
) OrderRouter {
	return &orderRouter{
		checkSession: checkSession,
		orderService: orderService,
	}
}

func (o *orderRouter) UploadOrderHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := o.checkSession(r)
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

	err = o.orderService.UploadOrder(login, order)
	if err != nil {
		if errors.Is(err, errors_custom.ErrOrderAlreadyUploadedBySameUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, errors_custom.ErrOrderAlreadyUploadedByAnotherUser) {
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
			return
		}
		http.Error(w, "Error uploading order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
func (o *orderRouter) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := o.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := o.orderService.GetOrders(login)
	if err != nil {
		if errors.Is(err, errors_custom.ErrNoContent) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.Error(w, "Error getting orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
func (o *orderRouter) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	login, ok := o.checkSession(r)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := o.orderService.GetBalance(login)
	if err != nil {
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(balance)
}
