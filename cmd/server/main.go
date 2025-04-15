package main

import (
	"context"
	"log"
	"pvz-service/config"
	"pvz-service/internal/database"
	"pvz-service/internal/logger"
)

func main() {
	ctx := context.Background()
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// init logger
	logger.InitLogger(cfg.LOG_LEVEL, cfg.MODE)

	// init db
	pool := database.ConnectDB(cfg, ctx)
	_ = pool
	log.Println("test docker")
}
