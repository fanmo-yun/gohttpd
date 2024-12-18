package main

import (
	"gohttpd/banner"
	web "gohttpd/internal"
	"gohttpd/logger"
	"gohttpd/utils"
	"log"
)

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("gohttpd: %v\n", err.Error())
	}
	cleanup := logger.NewLogger(config.Logger)
	defer cleanup()

	banner.ShowBanner()
	web.ServerRun(*config)
}
