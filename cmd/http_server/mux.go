package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k-akari/otel-example/internal/handler/httphandler"
	internal_middleware "github.com/k-akari/otel-example/internal/handler/httphandler/middleware"
)

func newMux(tsh *httphandler.TestService) http.Handler {
	middlewares := []func(http.Handler) http.Handler{
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60 * time.Second),
		internal_middleware.StartSpan,
	}

	mux := chi.NewRouter()
	mux.Use(middlewares...)

	mux.Get("/", tsh.PingGRPC)

	return mux
}
