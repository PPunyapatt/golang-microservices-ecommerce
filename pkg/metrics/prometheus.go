package metrics

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpMetrics struct {
	SuccessHttpRequests prometheus.Counter
	ErrorHttpRequests   prometheus.Counter

	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec

	AuthLoginRequests    prometheus.Counter
	AuthRegisterRequests prometheus.Counter

	PaymentWebhookRequests prometheus.Counter

	OrderPlaceRequests prometheus.Counter
}

type GrpcMetrics struct {
	SuccessRequests prometheus.Counter
	ErrorRequests   prometheus.Counter

	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec

	AuthLoginRequests    prometheus.Counter
	AuthRegisterRequests prometheus.Counter

	PaymentWebhookRequests prometheus.Counter

	OrderPlaceRequests prometheus.Counter
}

type Metrics struct {
	Http *HttpMetrics
	Grpc *GrpcMetrics
}

func NewMetrics() *Metrics {
	return &Metrics{
		&HttpMetrics{
			HttpRequestsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "http_requests_total",
					Help: "Total number of HTTP requests",
				},
				[]string{"status", "path", "method"},
			),
			SuccessHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
				Name: "success_http_requests_total",
				Help: "The total number of success http requests",
			}),
			ErrorHttpRequests: promauto.NewCounter(prometheus.CounterOpts{
				Name: "error_http_requests_total",
				Help: "The total number of error http requests",
			}),
			AuthLoginRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "auth_login_http_requests",
					Help: "Total number of auth login HTTP requests",
				}),
			AuthRegisterRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "auth_register_http_requests",
					Help: "Total number of auth register HTTP requests",
				}),
			PaymentWebhookRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "payment_webhook_http_requests",
					Help: "Total number of payment webhook HTTP requests",
				}),
			OrderPlaceRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "order_placement_http_requests",
					Help: "Total number of order placement HTTP requests",
				}),
			HttpRequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "http_request_duration_seconds",
					Help:    "Histogram of HTTP request latencies per endpoint",
					Buckets: prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, ...]
				},
				[]string{"method", "path", "status"},
			),
		},
		&GrpcMetrics{
			RequestsTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "grpc_requests_total",
					Help: "Total number of HTTP requests",
				},
				[]string{"service", "method", "status"},
			),
			RequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "grpc_request_duration_seconds",
					Help:    "Histogram of grpc request latencies per endpoint",
					Buckets: prometheus.DefBuckets, // [0.005, 0.01, 0.025, 0.05, ...]
				},
				[]string{"service", "method", "status"},
			),
			AuthLoginRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "auth_login_gRPC_requests_total",
					Help: "Total number of auth login gRPC requests",
				}),
			AuthRegisterRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "auth_register_gRPC_requests_total",
					Help: "Total number of auth register gRPC requests",
				}),
			PaymentWebhookRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "payment_webhook_gRPC_requests",
					Help: "Total number of payment webhook gRPC requests",
				}),
			OrderPlaceRequests: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "order_placement_gRPC_requests",
					Help: "Total number of order placement gRPC requests",
				}),
			SuccessRequests: promauto.NewCounter(prometheus.CounterOpts{
				Name: "success_grpc_requests_total",
				Help: "The total number of success grpc requests",
			}),
			ErrorRequests: promauto.NewCounter(prometheus.CounterOpts{
				Name: "error_grpc_requests_total",
				Help: "The total number of error grpc requests",
			}),
		},
	}
}

func (m *Metrics) PrometheusHttp(ctx context.Context, wg *sync.WaitGroup) {
	srv := &http.Server{
		Addr:    ":2112",
		Handler: promhttp.Handler(),
	}
	errCh := make(chan error, 1)

	go func() {
		log.Println("📊 prometheus http start")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("🛑 context done, shutting down server")

	case err := <-errCh:
		slog.Error("❌ prometheus http error", "err", err)
	}

	_ = srv.Shutdown(ctx)
	log.Println("🛑 Shutting prometheus http server...")
	wg.Done()
}
