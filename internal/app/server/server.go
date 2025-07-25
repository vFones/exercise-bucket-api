package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"bucket_organizer/internal/app/services"
	"bucket_organizer/internal/pkg/configs"
	"bucket_organizer/pkg/logger"
	"github.com/rs/cors"
)

type Server struct {
	router           *http.ServeMux
	httpServer       *http.Server
	services         *services.Services
	gracefulShutdown func(ctx context.Context) error // useful for shutting down clients e.g.: redis minio kafka
}

func NewServer(services *services.Services, shutdownFunc func(ctx context.Context) error) (*Server, error) {
	srv := &Server{
		router:           http.NewServeMux(),
		services:         services,
		gracefulShutdown: shutdownFunc,
	}
	srv.setupRoutes()
	return srv, nil
}

func (s *Server) Run(ctx context.Context) *http.Server {
	config := configs.Global().Server
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	s.httpServer = &http.Server{
		Addr:              ":" + strconv.Itoa(config.Port),
		Handler:           corsHandler.Handler(s.router),
		ReadTimeout:       time.Duration(config.Timeouts.Read) * time.Second,
		ReadHeaderTimeout: time.Duration(config.Timeouts.ReadHeaders) * time.Second,
		WriteTimeout:      time.Duration(config.Timeouts.Write) * time.Second,
		IdleTimeout:       time.Duration(config.Timeouts.Idle) * time.Second,
	}

	go func() {
		logger.Info(ctx, "Start serving http requests", logger.NewLogValue("port", config.Port))
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ctx, "HTTP server", err)
		}
	}()

	return s.httpServer
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.gracefulShutdown != nil {
		if err := s.gracefulShutdown(ctx); err != nil {
			logger.Error(ctx, "Shutdown server", err)
			return err
		}
	}
	return s.httpServer.Shutdown(ctx)
}
