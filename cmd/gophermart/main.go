package main

import (
	"log"

	"github.com/condratf/go-musthave-diploma-tpl/internal/config"
)

func main() {
	config.InitConfig()

	if err != nil {
		log.Fatal("server has crashed")
	}
}
