package orderupdater

import (
	"context"
	"time"

	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
)

type ordersService interface {
	GetPendingOrders(ctx context.Context) ([]models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string, accrual float64) error
}

type accrualService interface {
	GetOrderInfo(ctx context.Context, orderID string) (*models.AccrualResponse, error)
}

type orderUpdater struct {
	ordersService  ordersService
	accrualService accrualService
}

type OrderUpdater interface {
	StartOrderUpdater(ctx context.Context, interval time.Duration)
}
