package orderservice

import (
	"context"
	"regexp"

	"github.com/condratf/go-musthave-diploma-tpl/internal/errors_custom"
	"github.com/condratf/go-musthave-diploma-tpl/internal/models"
	"github.com/condratf/go-musthave-diploma-tpl/internal/utils"
)

func NewOrderService(
	ordersRepo orderRepository,
) OrderService {
	return &orderService{
		ordersRepo: ordersRepo,
	}
}

func (o *orderService) UploadOrder(login, order string) error {
	if !regexp.MustCompile(`^\d+$`).MatchString(order) {
		return errors_custom.ErrInvalidOrderNumber
	}

	if !utils.IsValidLuhn(order) {
		return errors_custom.ErrInvalidOrderNumber
	}

	err := o.ordersRepo.UploadOrder(login, order)
	if err != nil {
		return err
	}
	return nil
}

func (o *orderService) GetOrders(login string) ([]models.Order, error) {
	orders, err := o.ordersRepo.GetOrders(login)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, errors_custom.ErrNoContent
	}

	return orders, nil
}

func (o *orderService) GetPendingOrders(ctx context.Context) ([]models.Order, error) {
	orders, err := o.ordersRepo.GetPendingOrders(ctx)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return []models.Order{}, nil
	}

	return orders, nil
}

func (o *orderService) UpdateOrderStatus(ctx context.Context, orderNumber, status string, accrual float64) error {
	if !regexp.MustCompile(`^\d+$`).MatchString(orderNumber) {
		return errors_custom.ErrInvalidOrderNumber
	}

	if !utils.IsValidLuhn(orderNumber) {
		return errors_custom.ErrInvalidOrderNumber
	}

	err := o.ordersRepo.UpdateOrderStatus(ctx, orderNumber, status, accrual)
	if err != nil {
		return err
	}
	return nil
}

func (o *orderService) GetBalance(login string) (*float64, error) {
	balance, err := o.ordersRepo.GetBalance(login)
	if err != nil {
		return nil, err
	}

	if balance == nil {
		return nil, errors_custom.ErrNoContent
	}

	return &balance.Current, nil
}
