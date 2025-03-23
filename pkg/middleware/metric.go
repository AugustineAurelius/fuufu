package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type responseWriterWrapper struct {
	w          http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriterWrapper) Header() http.Header {
	return rw.w.Header()
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	size, err := rw.w.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.w.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func MetricMiddleware(meter metric.Meter, next http.Handler) http.Handler {
	requestCount, _ := meter.Int64Counter(
		"http.server.request.count",
	)

	durationHistogram, _ := meter.Int64Histogram(
		"http.server.duration",
		metric.WithDescription("Request duration in milliseconds"),
		metric.WithUnit("ms"),
	)

	responseSizeHistogram, _ := meter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("Response size in bytes"),
		metric.WithUnit("By"),
	)

	activeRequests, _ := meter.Int64UpDownCounter(
		"http.server.active_requests",
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := r.Context()

		baseLabels := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.route", r.URL.Path),
		}

		activeRequests.Add(ctx, 1, metric.WithAttributes(baseLabels...))
		defer activeRequests.Add(ctx, -1, metric.WithAttributes(baseLabels...))

		rw := &responseWriterWrapper{w: w}

		next.ServeHTTP(rw, r)

		baseLabels = append(baseLabels,
			attribute.Int("http.status_code", rw.statusCode),
		)

		requestCount.Add(ctx, 1, metric.WithAttributes(baseLabels...))

		duration := time.Since(startTime).Milliseconds()
		durationHistogram.Record(ctx, duration, metric.WithAttributes(baseLabels...))

		responseSizeHistogram.Record(ctx, int64(rw.size), metric.WithAttributes(baseLabels...))
	})
}
