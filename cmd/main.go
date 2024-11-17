package main

import (
	"gohttpd/banner"
	server "gohttpd/internal"
	"gohttpd/logger"
)

func main() {
	logger.InitLogger()
	banner.ShowBanner()
	server.ServerRun()
}
