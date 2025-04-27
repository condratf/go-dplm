package main

import (
	"fmt"
	"net/http"

	"github.com/condratf/go-musthave-diploma-tpl/internal/config"
	"github.com/condratf/go-musthave-diploma-tpl/internal/db"

	"github.com/condratf/go-musthave-diploma-tpl/internal/router"
	"github.com/condratf/go-musthave-diploma-tpl/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	config.InitConfig()
	db.InitDB()

	storage.NewPostgresStore(db.DB)

	r := chi.NewRouter()

	appRouter := router.Router(
		router.Config{
			RunAddress:           config.Config.RunAddress,
			DatabaseURI:          config.Config.DatabaseURI,
			AccrualSystemAddress: config.Config.AccrualSystemAddress,
		},
		db.DB,
	)

	r.Mount("/", appRouter)

	fmt.Println("Starting server on", config.Config.RunAddress)

	http.ListenAndServe(config.Config.RunAddress, r)
}
