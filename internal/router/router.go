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

	app := App{
		db:           db,
		config:       cnf,
		sessionStore: store,
		ordersRepo:   repository.NewOrdersRepository(db),
		userRepo:     repository.NewUserRepository(db),
	}

	r.Use(middleware.Logger)

	r.Post("/api/user/register", app.registerUserHandler)
	r.Post("/api/user/login", app.loginUserHandler)
	r.Post("/api/user/logout", app.authMiddleware(app.logoutUserHandler))
	r.Post("/api/user/orders", app.authMiddleware(app.uploadOrderHandler))
	r.Get("/api/user/orders", app.authMiddleware(app.getOrdersHandler))
	r.Get("/api/user/balance", app.authMiddleware(app.getBalanceHandler))
	r.Post("/api/user/balance/withdraw", app.authMiddleware(app.withdrawHandler))
	r.Get("/api/user/withdrawals", app.authMiddleware(app.getWithdrawalsHandler))

	return r
}
