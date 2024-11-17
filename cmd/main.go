package main

import (
	"gohttpd/banner"
	server "gohttpd/internal"
	"gohttpd/logger"
)

func main() {
	cleanup := logger.InitLogger()
	defer cleanup()
	banner.ShowBanner()
	server.ServerRun()
}
