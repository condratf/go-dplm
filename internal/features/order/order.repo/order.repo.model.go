package orderrepo

import (
	"context"
	"database/sql"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type OrderRepository interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error
	GetBalance(login string) (*models.BalanceResponse, error)
}

type orderRepository struct {
	db *sql.DB
}
