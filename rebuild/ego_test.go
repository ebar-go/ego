package rebuild

import "testing"

func TestRun(t *testing.T) {
	Run(ServerRunOptions{HttpAddr: ":8080"})
}