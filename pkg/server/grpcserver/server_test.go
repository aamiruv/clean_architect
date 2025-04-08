package grpcserver

import (
	"testing"
	"time"
)

func TestRunGrpcServer(t *testing.T) {
	server := New("localhost:12345", time.Second)
	go func() {
		if err := server.Run(); err != nil {
			t.Fail()
		}
	}()
	server.GracefulShutdown()
}
