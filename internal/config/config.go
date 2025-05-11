package config

import (
	"flag"
	"os"
)

type config struct {
	RunAddress           string
	AccrualSystemAddress string
	DatabaseURI          string
}

var Config = config{
	RunAddress:           "localhost:8080",
	AccrualSystemAddress: "http://accrual-system",
	// DatabaseURI:          "",
	DatabaseURI: "postgres://user:password@localhost/gophermart?sslmode=disable",
}

func InitConfig() {
	runAddr := flag.String("a", "", "HTTP server address")
	accrualAddr := flag.String("r", "", "accrual address")
	databaseURI := flag.String("d", "", "Database URI")

	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		Config.RunAddress = envAddr
	} else if *runAddr != "" {
		Config.RunAddress = *runAddr
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		Config.DatabaseURI = envDatabaseURI
	} else if *databaseURI != "" {
		Config.DatabaseURI = *databaseURI
	}

	if envAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddr != "" {
		Config.AccrualSystemAddress = envAccrualAddr
	} else if *accrualAddr != "" {
		Config.AccrualSystemAddress = *accrualAddr
	}
}
