package main

import (
	"context"
	"pvz-service/config"
	"pvz-service/internal/database"
	"pvz-service/internal/logger"

	"github.com/labstack/echo"
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
	// init echo
	e := echo.New()
	e.Logger.Fatal(e.Start(":8080"))

	//graseful shutdown
}
