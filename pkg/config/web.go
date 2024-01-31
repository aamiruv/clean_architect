package config

import "fmt"

type web struct {
	BindingIpAddress string `json:"bindingIpAddress"`
	Port             uint   `json:"port"`
}

func (w web) GetAddress() string {
	return fmt.Sprintf("%s:%d", w.BindingIpAddress, w.Port)
}
