package app

import (
	"net"
	"payment/v1/proto/payment"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func StartgRPCServer(paymentServiceRPC payment.PaymentServiceServer) {
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	listener, err := net.Listen("tcp", ":1029")
	if err != nil {
		panic(err)
	}

	payment.RegisterPaymentServiceServer(s, paymentServiceRPC)

	// ✅ Register health check service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)

	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
}
