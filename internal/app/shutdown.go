package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bucket_organizer/internal/app/server"
	"bucket_organizer/pkg/logger"
)

func GracefulShutdown(ctx context.Context, timeout time.Duration, server *server.Server, c chan os.Signal) {
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info(ctx, "Shutting down server.", logger.NewLogValue("timeout", timeout.String()))
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer func() {
		_ = logger.Sync()
		cancel()

	}()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error(ctx, "server graceful shutdown failed.", err)
	}
}
