package main

import (
	"context"
	"os"
	"time"

	"bucket_organizer/internal/app"
	"bucket_organizer/internal/app/server/dependency"
	"bucket_organizer/internal/pkg/configs"
	"bucket_organizer/pkg/logger"
)

func Run(ctx context.Context) error {
	if err := configs.LoadConfig(); err != nil {
		return err
	}
	config := configs.Global()
	logger.InitLogger(configs.IsDevelopment())

	logger.Debug(ctx, "configuration", logger.NewLogValue("config", config))

	srv, err := dependency.Inject(ctx)
	if err != nil {
		return err
	}

	_ = srv.Run(ctx)
	timeout := time.Duration(config.Server.Timeouts.Shutdown) * time.Second
	c := make(chan os.Signal, 1)
	app.GracefulShutdown(ctx, timeout, srv, c)
	return nil
}
