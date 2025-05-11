package app

import (
	"context"
	"database/sql"
	"fmt"

	"net/http"
	"time"

	"github.com/condratf/go-musthave-diploma-tpl/internal/config"
	accrual_service "github.com/condratf/go-musthave-diploma-tpl/internal/features/accrual/accrual.service"
	order_updater "github.com/condratf/go-musthave-diploma-tpl/internal/features/order-updater"
	order_repo "github.com/condratf/go-musthave-diploma-tpl/internal/features/order/order.repo"
	order_router "github.com/condratf/go-musthave-diploma-tpl/internal/features/order/order.router"
	order_service "github.com/condratf/go-musthave-diploma-tpl/internal/features/order/order.service"
	user_repo "github.com/condratf/go-musthave-diploma-tpl/internal/features/user/user.repo"
	user_router "github.com/condratf/go-musthave-diploma-tpl/internal/features/user/user.router"
	user_service "github.com/condratf/go-musthave-diploma-tpl/internal/features/user/user.service"
	withdraw_repo "github.com/condratf/go-musthave-diploma-tpl/internal/features/withdraw/withdraw.repo"
	withdraw_router "github.com/condratf/go-musthave-diploma-tpl/internal/features/withdraw/withdraw.router"
	withdraw_service "github.com/condratf/go-musthave-diploma-tpl/internal/features/withdraw/withdraw.service"
	"github.com/condratf/go-musthave-diploma-tpl/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

// todo:
const secretKey = "[secret-key]"

func RunApp(
	cnf struct {
		RunAddress           string
		DatabaseURI          string
		AccrualSystemAddress string
	},
	db *sql.DB,
	ctx context.Context,
) error {
	orderRepo := order_repo.NewOrderRepository(db)
	userRepo := user_repo.NewUserRepository(db)
	withdrawRepo := withdraw_repo.NewWithdrawRepository(db)

	userService := user_service.NewUserService(userRepo)
	orderService := order_service.NewOrderService(orderRepo)
	withdrawService := withdraw_service.NewWithdrawService(withdrawRepo)

	accrualService := accrual_service.NewAccrualClient(cnf.AccrualSystemAddress, &http.Client{Timeout: 11 * time.Second})
	updater := order_updater.NewOrderUpdater(orderService, accrualService)
	updater.StartOrderUpdater(ctx, 10*time.Second)

	sessionsStore := sessions.NewCookieStore([]byte(secretKey))

	userRouter := user_router.NewUserRouter(sessionsStore, userService)
	orderRouter := order_router.NewOrderRouter(userRouter.CheckSession, orderService)
	withdrawRouter := withdraw_router.NewWithdrawRouter(userRouter.CheckSession, withdrawService)

	r := chi.NewRouter()

	r.Post("/api/user/register", userRouter.RegisterUserHandler)
	r.Post("/api/user/login", userRouter.LoginUserHandler)
	r.Post("/api/user/logout", utils.AuthMiddleware(sessionsStore, userRouter.LogoutUserHandler))
	r.Post("/api/user/orders", utils.AuthMiddleware(sessionsStore, orderRouter.UploadOrderHandler))
	r.Get("/api/user/orders", utils.AuthMiddleware(sessionsStore, orderRouter.GetOrdersHandler))
	r.Get("/api/user/balance", utils.AuthMiddleware(sessionsStore, orderRouter.GetBalanceHandler))
	r.Post("/api/user/balance/withdraw", utils.AuthMiddleware(sessionsStore, withdrawRouter.WithdrawHandler))
	r.Get("/api/user/withdrawals", utils.AuthMiddleware(sessionsStore, withdrawRouter.GetWithdrawalsHandler))

	fmt.Println("Starting server on", config.Config.RunAddress)
	return http.ListenAndServe(config.Config.RunAddress, r)
}
