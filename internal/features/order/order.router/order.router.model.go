package orderrouter

import (
	"context"
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type OrderRouter interface {
	UploadOrderHandler(w http.ResponseWriter, r *http.Request)
	GetOrdersHandler(w http.ResponseWriter, r *http.Request)
	GetBalanceHandler(w http.ResponseWriter, r *http.Request)
}

type orderService interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetBalance(login string) (*float64, error)
}

type accrualService interface {
	RegisterOrder(ctx context.Context, orderNumber string) error
}

type orderRouter struct {
	checkSession   func(r *http.Request) (string, bool)
	orderService   orderService
	accrualService accrualService
}
