package example

import (
	"github.com/ebar-go/ego"
	"testing"
)

func TestAggregatorWithGrpcServer(t *testing.T) {
	aggregator := ego.NewAggregatorServer()

	grpcServer := ego.NewGRPCServer(":8081")

	aggregator.WithServer(grpcServer)

	aggregator.Run()
}
