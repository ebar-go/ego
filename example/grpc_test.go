package example

import (
	"github.com/ebar-go/ego"
	"testing"
)

func TestAggregatorWithGrpcServer(t *testing.T) {
	aggregator := ego.NewAggregatorServer()

	httpServer := ego.NewGRPCServer(":8081")

	aggregator.WithServer(httpServer)

	aggregator.Run()
}
