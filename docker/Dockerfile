# Step 1: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# คัดลอก go.work, go.mod, go.sum เพื่อ resolve dependency
COPY go.work go.work.sum ./
COPY pkg ./pkg
COPY services/auth ./services/auth

WORKDIR /app/services/auth/v1/cmd

# ดึง dependency
RUN go mod tidy

# Build binary
RUN go build -o /app/bin/auth-service main.go

# Step 2: Run stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/bin/auth-service .

CMD ["./auth-service"]
