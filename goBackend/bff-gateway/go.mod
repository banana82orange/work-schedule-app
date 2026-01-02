module github.com/portfolio/bff-gateway

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/portfolio/proto v0.0.0
	github.com/portfolio/shared v0.0.0
	google.golang.org/grpc v1.64.0
	google.golang.org/protobuf v1.34.0
)

require github.com/joho/godotenv v1.5.1 // indirect

replace github.com/portfolio/shared => ../shared

replace github.com/portfolio/proto => ../proto
