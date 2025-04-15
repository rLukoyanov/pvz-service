package main

import (
	"log"
	"pvz-service/config"
	"pvz-service/internal/logger"
)

func main() {
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	_ = cfg

	// init logger
	logger.InitLogger(cfg.LOG_LEVEL, cfg.MODE)

	// init db
	log.Println("test docker")
}
