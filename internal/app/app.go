package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/condratf/go-musthave-diploma-tpl/internal/config"
	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/condratf/go-musthave-diploma-tpl/internal/router"
	"github.com/condratf/go-musthave-diploma-tpl/internal/services/accrual"
	"github.com/go-chi/chi/v5"
)

type App struct {
	ordersRepo     repository.OrdersRepository
	accrualService accrual.AccrualService
}

func RunApp(
	cnf struct {
		RunAddress           string
		DatabaseURI          string
		AccrualSystemAddress string
	},
	db *sql.DB,
	ctx context.Context,
) error {
	app := &App{
		ordersRepo:     repository.NewOrdersRepository(db),
		accrualService: accrual.NewAccrualClient(cnf.AccrualSystemAddress, &http.Client{Timeout: 10 * time.Second}),
	}

	app.startOrderUpdater(ctx, 30*time.Second)

	r := chi.NewRouter()
	appRouter := router.Router(
		router.Config{
			RunAddress:           config.Config.RunAddress,
			DatabaseURI:          config.Config.DatabaseURI,
			AccrualSystemAddress: config.Config.AccrualSystemAddress,
		},
		db,
	)

	r.Mount("/", appRouter)
	fmt.Println("Starting server on", config.Config.RunAddress)
	return http.ListenAndServe(config.Config.RunAddress, r)
}

func (app *App) startOrderUpdater(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				app.updateOrders(ctx)
			case <-ctx.Done():
				log.Println("Order updater stopped")
				return
			}
		}
	}()
}

func (app *App) updateOrders(ctx context.Context) {
	orders, err := app.ordersRepo.GetPendingOrders(ctx)
	if err != nil {
		log.Println("Error fetching pending orders:", err)
		return
	}

	for _, order := range orders {
		resp, err := app.accrualService.GetOrderInfo(ctx, order.Order)
		if err != nil {
			log.Println("Error fetching accrual info for order", order.Order, ":", err)
			continue
		}

		if resp.Status == "INVALID" || resp.Status == "PROCESSED" {
			if err := app.ordersRepo.UpdateOrderStatus(ctx, order.Order, resp.Status, resp.Accrual); err != nil {
				log.Println("Error updating order status:", err)
			}
		} else {
			if err := app.ordersRepo.UpdateOrderStatus(ctx, order.Order, resp.Status, 0); err != nil {
				log.Println("Error updating order status:", err)
			}
		}
	}
}
