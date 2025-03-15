package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware(tracer trace.Tracer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "middleware", trace.WithAttributes(
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.String("proto", r.Proto),
			attribute.String("remote_addr", r.RemoteAddr),
		))
		defer span.End()

		next.ServeHTTP(w, r.Clone(ctx))
	})
}
