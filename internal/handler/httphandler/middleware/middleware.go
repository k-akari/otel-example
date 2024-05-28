package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type Middleware struct {
	tracer trace.Tracer
}

func New(tracer trace.Tracer) *Middleware {
	return &Middleware{tracer: tracer}
}

func (m *Middleware) StartSpan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := m.tracer.Start(r.Context(), fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path))
		defer span.End()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
