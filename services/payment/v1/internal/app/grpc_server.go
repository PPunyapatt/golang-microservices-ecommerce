package app

import (
	"context"
	"log/slog"
	"net"
	"package/interceptor"
	"package/metrics"
	"payment/v1/proto/payment"
	"sync"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartgRPCServer(ctx context.Context, paymentServiceRPC payment.PaymentServiceServer, wg *sync.WaitGroup, pm *metrics.Metrics) {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor(pm)),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1029")
	if err != nil {
		panic(err)
	}

	payment.RegisterPaymentServiceServer(s, paymentServiceRPC)

	slog.Info("🚀 gRPC server started on :1029")
	go func() {
		if err := s.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			slog.Error("gRPC server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("🛑 Shutting down gRPC server...")
	wg.Done()
	s.GracefulStop()
}
