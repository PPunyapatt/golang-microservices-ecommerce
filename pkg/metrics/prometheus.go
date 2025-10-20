package metrics

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	SuccessHttpRequests prometheus.Counter
	ErrorHttpRequests   prometheus.Counter
	AuthLoginRequests   prometheus.Counter
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	return &Metrics{
		HttpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"status", "path", "method"},
		),
		AuthLoginRequests: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "auth_login_http_requests_total",
				Help: "Total number of auth login HTTP requests",
			}),
		HttpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of HTTP request latencies per endpoint",
				Buckets: prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, ...]
			},
			[]string{"method", "path", "status"},
		),
	}
}

func (m *Metrics) PrometheusHttp(ctx context.Context, wg *sync.WaitGroup) {
	srv := &http.Server{
		Addr:    ":2112",
		Handler: promhttp.Handler(),
	}

	go func() {
		log.Println("📊 prometheus http start")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ prometheus http error: %v\n", err)
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
	log.Println("🛑 Shutting prometheus http server...")
	wg.Done()
}
