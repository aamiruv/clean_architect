package config

import "fmt"

type grpc struct {
	BindingIpAddress string `json:"bindingIpAddress"`
	Port             uint   `json:"port"`
}

func (g grpc) Address() string {
	return fmt.Sprintf("%s:%d", g.BindingIpAddress, g.Port)
}
