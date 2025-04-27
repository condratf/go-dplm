package router

import (
	"database/sql"

	"github.com/condratf/go-musthave-diploma-tpl/internal/repository"
	"github.com/gorilla/sessions"
)

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
}

type App struct {
	db           *sql.DB
	sessionStore *sessions.CookieStore
	userRepo     repository.UserRepository
	ordersRepo   repository.OrdersRepository
	config       Config
}
