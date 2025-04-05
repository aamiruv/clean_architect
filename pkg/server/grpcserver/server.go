// Pacakge grpcserver provides grpc server construction.
package grpcserver

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	grpcServer      *grpc.Server
	bindingAddress  string
	shutdownTimeout time.Duration
}

func New(bindingAddress string, shutdownTimeout time.Duration, options ...grpc.ServerOption) server {
	s := grpc.NewServer(options...)
	return server{
		grpcServer:      s,
		bindingAddress:  bindingAddress,
		shutdownTimeout: shutdownTimeout,
	}
}

func (s server) Run() error {
	ls, err := net.Listen("tcp", s.bindingAddress)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(ls)
}

func (s server) GracefulShutdown() {
	s.grpcServer.GracefulStop()
	time.AfterFunc(s.shutdownTimeout, s.grpcServer.Stop)
}

func (s server) Server() *grpc.Server {
	return s.grpcServer
}
