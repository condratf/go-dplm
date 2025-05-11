package orderservice

import (
	"context"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type OrderService interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error
	GetBalance(login string) (*float64, error)
}

type orderRepository interface {
	UploadOrder(login, order string) error
	GetOrders(login string) ([]models.Order, error)
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error
	GetBalance(login string) (*models.BalanceResponse, error)
}

type orderService struct {
	ordersRepo orderRepository
}
