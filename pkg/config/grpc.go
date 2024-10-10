package config

import "fmt"

type grpc struct {
	BindingIpAddress  string `json:"bindingIpAddress"`
	Port              uint   `json:"port"`
	MaxReceiveMsgSize int    `json:"maxReceiveMsgSize"`
	ReadBufferSize    int    `json:"readBufferSize"`
	HasReflection     bool   `json:"hasReflection"`
}

func (g grpc) Address() string {
	return fmt.Sprintf("%s:%d", g.BindingIpAddress, g.Port)
}
