package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mustaphalimar/gopanel/internal/database"
	"github.com/mustaphalimar/gopanel/internal/router"
	"github.com/mustaphalimar/gopanel/internal/server"
	"github.com/mustaphalimar/gopanel/pkg/config"
	"github.com/mustaphalimar/gopanel/pkg/logger"
)

type App struct {
	Cfg *config.Config
	Log *logger.Logger
	DB  *database.Database
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := logger.New(&cfg.Logger)

	db, err := database.New(&cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	return &App{
		Cfg: cfg,
		Log: log,
		DB:  db,
	}, nil
}

func (a *App) RunAPI() error {
	a.Log.Info("starting gopanel API",
		"env", a.Cfg.Server.Env,
		"port", a.Cfg.Server.Port,
	)

	// run db migrations
	if err := database.Migrate(context.Background(), a.Cfg, a.Log); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	srv, err := server.New(a.Cfg, a.DB, a.Log)
	if err != nil {
		return fmt.Errorf("init server: %w", err)
	}

	r := router.New(a.Log)
	srv.SetupHTTPServer(r)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.Start()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("http server: %w", err)
	case sig := <-shutdown:
		a.Log.Info("shutdown signal received", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 16*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			a.Log.Error("graceful shutdown failed", "error", err)
		}
	}
	return nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
