package server

import (
	"net/http"

	"github.com/mustaphalimar/gopanel/internal/database"
	"github.com/mustaphalimar/gopanel/pkg/config"
	"github.com/mustaphalimar/gopanel/pkg/logger"
)

type Server struct {
	Config   *config.Config
	Logger   *logger.Logger
	Database *database.PgDB

	httpServer *http.Server
}
