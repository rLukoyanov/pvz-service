package main

import "pvz-service/config"

func main() {
	// init config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	_ = cfg

	// init logger

	// init db
}
