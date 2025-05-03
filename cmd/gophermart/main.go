package main

import (
	"context"
	"log"

	"github.com/condratf/go-musthave-diploma-tpl/internal/app"
	"github.com/condratf/go-musthave-diploma-tpl/internal/config"
	"github.com/condratf/go-musthave-diploma-tpl/internal/db"

	"github.com/condratf/go-musthave-diploma-tpl/internal/storage"
)

func main() {
	config.InitConfig()
	db.InitDB()
	storage.NewPostgresStore(db.DB)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	error := app.RunApp(struct {
		RunAddress           string
		DatabaseURI          string
		AccrualSystemAddress string
	}{
		RunAddress:           config.Config.RunAddress,
		DatabaseURI:          config.Config.DatabaseURI,
		AccrualSystemAddress: config.Config.AccrualSystemAddress,
	}, db.DB, ctx)

	if error != nil {
		log.Fatal("server has crashed")
	}
}
