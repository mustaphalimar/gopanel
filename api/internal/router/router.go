package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/mustaphalimar/gopanel/pkg/logger"
)

type Router struct {
	router *chi.Mux
	logger *logger.Logger
}

func New(logger *logger.Logger) *Router {
	r := &Router{
		router: chi.NewRouter(),
		logger: logger,
	}

	r.setupMiddlewares()
	r.setupRoutes()

	return r
}

func (r *Router) setupMiddlewares() {
	r.router.Use(chimiddleware.RequestID)
	// TODO add an alternative to RealIP
	// r.router.Use(chimiddleware.RealIP)

}

func (r *Router) setupRoutes() {
	r.router.Route("/api", func(r chi.Router) {

	})
}

func (r *Router) ServerHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
