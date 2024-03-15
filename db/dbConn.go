package db

import (
	"database/sql"
	"github.com/aalperen0/dailytask/internal/database"
	"log"
	"os"
)

type ApiConfig struct {
	DB *database.Queries
}

func InitAPIConfig() *ApiConfig {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Failed to get url or url is wrong")
	}

	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect database", err)
	}

	apiCfg := ApiConfig{
		DB: database.New(connection),
	}

	return &apiCfg
}
