package main

import (
	"gohttpd/banner"
	"gohttpd/logger"
)

func main() {
	cleanup := logger.InitLogger()
	defer cleanup()
	banner.ShowBanner()
	ParseCommandAndRun()
}
