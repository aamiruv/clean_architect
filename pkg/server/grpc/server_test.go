package grpc

import (
	"testing"
	"time"
)

func TestRunGrpcServer(t *testing.T) {
	server := NewServer("localhost:12345")
	go func() {
		if err := server.Run(); err != nil {
			t.Fail()
		}
	}()
	server.GracefulShutdown(time.Second)
}
