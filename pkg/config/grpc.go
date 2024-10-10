package config

import "fmt"

type grpc struct {
	bindingIpAddress  string
	port              uint
	maxReceiveMsgSize int
	readBufferSize    int
	hasReflection     bool
}

func (g grpc) Address() string {
	return fmt.Sprintf("%s:%d", g.bindingIpAddress, g.port)
}

func (g grpc) MaxReceiveMsgSize() int {
	return g.maxReceiveMsgSize
}

func (g grpc) ReadBufferSize() int {
	return g.readBufferSize
}

func (g grpc) HasReflection() bool {
	return g.hasReflection
}
