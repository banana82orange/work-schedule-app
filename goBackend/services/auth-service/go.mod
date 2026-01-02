module github.com/portfolio/auth-service

go 1.21

require (
	github.com/portfolio/proto v0.0.0
	github.com/portfolio/shared v0.0.0
	golang.org/x/crypto v0.17.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.0
)

replace github.com/portfolio/shared => ../../shared

replace github.com/portfolio/proto => ../../proto
