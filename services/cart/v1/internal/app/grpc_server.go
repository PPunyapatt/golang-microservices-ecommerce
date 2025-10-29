package app

import (
	"cart/v1/proto/cart"
	"context"
	"log/slog"
	"net"
	"sync"

	"package/interceptor"
	"package/metrics"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartgRPCServer(ctx context.Context, cartServerRPC cart.CartServiceServer, wg *sync.WaitGroup, pm *metrics.Metrics) {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor(pm)),
	)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	listener, err := net.Listen("tcp", ":1027")
	if err != nil {
		panic(err)
	}

	cart.RegisterCartServiceServer(s, cartServerRPC)

	slog.Info("🚀 gRPC server started on :1027")
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
