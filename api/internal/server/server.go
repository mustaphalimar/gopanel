package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mustaphalimar/gopanel/internal/database"
	"github.com/mustaphalimar/gopanel/pkg/config"
	"github.com/mustaphalimar/gopanel/pkg/logger"
)

type Server struct {
	Config   *config.Config
	Logger   *logger.Logger
	Database *database.Database

	httpServer *http.Server
}

func New(cfg *config.Config, db *database.Database, logger *logger.Logger) (*Server, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	server := &Server{
		Config:   cfg,
		Logger:   logger,
		Database: db,
	}
	return server, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port),
		Handler:      handler,
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
		IdleTimeout:  s.Config.Server.IdleTimeout,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}
	s.Logger.Info("Starting HTTP server", "address", fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.Logger.Info("Shutting down HTTP server", "address", fmt.Sprintf("%s:%s", s.Config.Server.Host, s.Config.Server.Port))
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}
