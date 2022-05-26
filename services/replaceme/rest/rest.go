package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iconimpact/replaceme/internal/mongodb"
	"github.com/iconimpact/replaceme/internal/redisdb"
	"github.com/iconimpact/replaceme/internal/sqldb"
	"go.uber.org/zap"
)

type REST struct {
	Router *chi.Mux
	srv    *http.Server

	DB    *sqldb.Client
	Mongo *mongodb.Client
	Redis *redisdb.Client

	logger *zap.SugaredLogger
}

func New(DB *sqldb.Client, Mongo *mongodb.Client, redis *redisdb.Client, port string, logger *zap.SugaredLogger) *REST {
	rest := &REST{
		DB:     DB,
		Mongo:  Mongo,
		Redis:  redis,
		logger: logger,
	}

	mux := chi.NewRouter()

	// Add middlewares
	mux.Use(rest.LogRequestMiddleware)
	mux.Use(addTimeContextMiddleware) // used for request-time and action-time headers
	//r.Use(timeTrackingMiddleware)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(180 * time.Second))

	rest.Router = mux
	rest.AddRoutes()

	rest.srv = &http.Server{Addr: ":" + port, Handler: rest.Router}

	return rest
}

func (rest *REST) Start() {
	rest.logger.Infow("Starting REST service", "addr", rest.srv.Addr)

	if err := rest.srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		rest.logger.Fatalf("healthcheck server error: %v", err)
	}
}

func (rest *REST) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := rest.srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		rest.logger.Errorf("healthcheck server shutdown error: %v", err)
	}
}
