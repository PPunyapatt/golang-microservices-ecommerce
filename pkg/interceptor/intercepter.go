package interceptor

import (
	"context"
	"fmt"
	"log/slog"
	"package/metrics"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor(pm *metrics.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		spanCtx := trace.SpanContextFromContext(ctx)
		traceID := spanCtx.TraceID().String()

		parts := strings.Split(info.FullMethod, "/")
		service := "unknown"
		method := "unknown"
		if len(parts) == 3 {
			service = parts[1] // "auth.AuthService"
			method = parts[2]  // "Login"
		}

		resp, err := handler(ctx, req)

		level := slog.LevelInfo
		status, errMsg, color := "OK", "-", "🔵"

		if err != nil {
			status, errMsg, color, level = "ERROR", err.Error(), "🔴", slog.LevelError
			pm.Grpc.ErrorRequests.Inc()
		} else {
			pm.Grpc.SuccessRequests.Inc()
		}

		duration := time.Since(start).Seconds()

		slog.Log(
			ctx,
			level,
			fmt.Sprintf("%s %s | %s | %.4fms | %s | %s | %s | %s",
				color,
				time.Now().Format("15:04:05"),
				status,
				duration,
				service,
				method,
				traceID,
				errMsg,
			))

		pm.Grpc.RequestsTotal.WithLabelValues(service, method, status).Inc()
		pm.Grpc.RequestDuration.WithLabelValues(service, method, status).Observe(duration)

		return resp, err
	}
}
