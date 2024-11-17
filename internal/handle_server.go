package internal

import (
	"context"
	"fmt"
	"gohttpd/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func ServerRun(c utils.Config) {
	address := fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
	gohttp := http.Server{
		Addr:    address,
		Handler: HandleRouter(c),
	}

	go func() {
		zap.L().Info("Server Running", zap.String("http", address))
		err := gohttp.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("gohttpd: Listen And Serve Fatal", zap.Error(err))
		}
	}()

	ShutdownServer(&gohttp)
}

func ShutdownServer(s *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	zap.L().Warn("gohttpd: Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		zap.L().Fatal("gohttpd: Server forced to shutdown", zap.Error(err))
	}

	zap.L().Warn("gohttpd: Server exiting")
}
