// Pacakge grpcserver provides grpc server construction.
package grpcserver

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	*grpc.Server
	bindingAddress  string
	shutdownTimeout time.Duration
}

func New(bindingAddress string, shutdownTimeout time.Duration, options ...grpc.ServerOption) server {
	s := grpc.NewServer(options...)
	return server{
		s,
		bindingAddress,
		shutdownTimeout,
	}
}

func (s *server) Run() error {
	ls, err := net.Listen("tcp", s.bindingAddress)
	if err != nil {
		return err
	}
	return s.Serve(ls)
}

func (s *server) GracefulShutdown() {
	s.GracefulStop()
	time.AfterFunc(s.shutdownTimeout, s.Stop)
}
