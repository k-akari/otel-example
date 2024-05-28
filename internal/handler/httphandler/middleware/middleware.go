package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

func StartSpan(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const instrumentationName = "github.com/k-akari/otel-example/internal/handler/httphandler/middleware"
		ctx, span := otel.Tracer(instrumentationName).Start(r.Context(), fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path))
		defer span.End()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
