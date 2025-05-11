package orderupdater

import (
	"context"
	"log"
	"time"
)

func NewOrderUpdater(
	ordersService ordersService,
	accrualService accrualService,
) OrderUpdater {
	return &orderUpdater{
		ordersService:  ordersService,
		accrualService: accrualService,
	}
}

func (o *orderUpdater) StartOrderUpdater(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				o.updateOrders(ctx)
			case <-ctx.Done():
				log.Println("Order updater stopped")
				return
			}
		}
	}()
}

func (o *orderUpdater) updateOrders(ctx context.Context) {
	orders, err := o.ordersService.GetPendingOrders(ctx)
	if err != nil {
		log.Println("Error fetching pending orders:", err)
		return
	}

	for _, order := range orders {
		resp, err := o.accrualService.GetOrderInfo(ctx, order.Order)
		if err != nil {
			log.Println("Error fetching accrual info for order", order.Order, ":", err)
			continue
		}

		if resp.Status == "INVALID" || resp.Status == "PROCESSED" {
			if err := o.ordersService.UpdateOrderStatus(ctx, order.Order, resp.Status, resp.Accrual); err != nil {
				log.Println("Error updating order status:", err)
			}
		} else {
			if err := o.ordersService.UpdateOrderStatus(ctx, order.Order, resp.Status, 0); err != nil {
				log.Println("Error updating order status:", err)
			}
		}
	}
}
