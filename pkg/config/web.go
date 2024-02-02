package config

import (
	"fmt"
	"time"
)

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

func (w web) GetReadTimeOut() time.Duration {
	return time.Duration(w.ReadTimeOutInSec) * time.Second
}

func (w web) GetIdleTimeout() time.Duration {
	return time.Duration(w.IdleTimeoutInSec) * time.Second
}

func (w web) GetWriteTimeout() time.Duration {
	return time.Duration(w.WriteTimeoutInSec) * time.Second
}

func (w web) GetReadHeaderTimeout() time.Duration {
	return time.Duration(w.ReadHeaderTimeoutInSec) * time.Second
}
