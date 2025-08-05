package helper

import (
	"context"
	"gateway/v1/internal/constant"
	"gateway/v1/proto/auth"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func ConnectGRPC(address string) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server at %s: %v", address, err)
	}

	log.Println("Connected to gRPC server at", address)
	healthClient := grpc_health_v1.NewHealthClient(conn)
	_, err = healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		conn.Close()
		log.Fatalf("Health check failed for %s: %v", address, err.Error())
	}
	return conn
}

func NewClientsGRPC() *constant.Clients {
	userConn := ConnectGRPC(os.Getenv("auth"))
	// cartConn := ConnectGRPC(os.Getenv("cart"))
	// inventoryConn := ConnectGRPC(os.Getenv("inventory"))

	return &constant.Clients{
		// CartClient: cart.NewCartServiceClient(cartConn),
		AuthClient: auth.NewAuthServiceClient(userConn),
		// InventoryClient: Inventory.NewInventoryServiceClient(inventoryConn),
	}
}
