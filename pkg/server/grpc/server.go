package grpc

import (
	"google.golang.org/grpc"
	"net"
	"time"
)

type server struct {
	grpcServer     *grpc.Server
	bindingAddress string
}

func NewServer(bindingAddress string, options ...grpc.ServerOption) server {
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

func (s server) GetServer() *grpc.Server {
	return s.grpcServer
}
