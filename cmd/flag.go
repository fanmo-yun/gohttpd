package main

import (
	"flag"
	web "gohttpd/internal"
	"gohttpd/server"
	"gohttpd/utils"

	"go.uber.org/zap"
)

func ParseCommandAndRun() {
	var ServerType string
	flag.StringVar(&ServerType, "type", "web", "gohttpd server mode (etcd or web)")

	flag.Parse()

	if ServerType == "" {
		zap.L().Fatal("gohttpd: -type parameter is required.")
	}

	config := utils.LoadConfig()

	switch ServerType {
	case "etcd":
		server.EtcdServiceRegister()
	case "web":
		web.ServerRun(*config)
	default:
		zap.L().Fatal("gohttpd: Unknown server type", zap.String("type", ServerType))
	}
}
