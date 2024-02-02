package config

import "fmt"

type web struct {
	BindingIpAddress       string `json:"bindingIpAddress"`
	Port                   uint   `json:"port"`
	ReadTimeOutInSec       uint   `json:"readTimeOutInSec"`
	IdleTimeoutInSec       uint   `json:"idleTimeoutInSec"`
	WriteTimeoutInSec      uint   `json:"writeTimeoutInSec"`
	ReadHeaderTimeoutInSec uint   `json:"readHeaderTimeoutInSec"`
}

func (w web) GetAddress() string {
	return fmt.Sprintf("%s:%d", w.BindingIpAddress, w.Port)
}
