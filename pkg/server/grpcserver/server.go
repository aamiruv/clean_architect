// Pacakge grpcserver provides grpc server construction.
package grpcserver

import (
	"net"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	grpcServer     *grpc.Server
	bindingAddress string
}

func New(bindingAddress string, options ...grpc.ServerOption) server {
	s := grpc.NewServer(options...)
	return server{grpcServer: s, bindingAddress: bindingAddress}
}

func (s server) Run() error {
	ls, err := net.Listen("tcp", s.bindingAddress)
	if err != nil {
		return err
	}
	return s.grpcServer.Serve(ls)
}

func (s server) GracefulShutdown(deadline time.Duration) {
	s.grpcServer.GracefulStop()
	go func() {
		time.Sleep(deadline)
		s.grpcServer.Stop()
	}()
}

func (s server) Server() *grpc.Server {
	return s.grpcServer
}
