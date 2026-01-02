# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git

# Copy proto module
COPY proto/ ./proto/

# Copy shared module
COPY shared/ ./shared/

# Copy service module
COPY services/auth-service/ ./services/auth-service/

# Generate go.sum and download dependencies
WORKDIR /build/proto
RUN go mod tidy

WORKDIR /build/shared
RUN go mod tidy

WORKDIR /build/services/auth-service
RUN go mod tidy
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /auth-service .
EXPOSE 50051
CMD ["./auth-service"]
