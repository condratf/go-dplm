package router

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func Router(
	cnf Config,
	db *sql.DB,
) http.Handler {
	r := chi.NewRouter()

	appRouter := AppRouter{
		db:           db,
		config:       cnf,
		sessionStore: store,
		ordersRepo:   repository.NewOrdersRepository(db),
		userRepo:     repository.NewUserRepository(db),
	}

	r.Use(middleware.Logger)

	r.Post("/api/user/register", appRouter.registerUserHandler)
	r.Post("/api/user/login", appRouter.loginUserHandler)
	r.Post("/api/user/logout", appRouter.authMiddleware(appRouter.logoutUserHandler))
	r.Post("/api/user/orders", appRouter.authMiddleware(appRouter.uploadOrderHandler))
	r.Get("/api/user/orders", appRouter.authMiddleware(appRouter.getOrdersHandler))
	r.Get("/api/user/balance", appRouter.authMiddleware(appRouter.getBalanceHandler))
	r.Post("/api/user/balance/withdraw", appRouter.authMiddleware(appRouter.withdrawHandler))
	r.Get("/api/user/withdrawals", appRouter.authMiddleware(appRouter.getWithdrawalsHandler))

	return r
}
