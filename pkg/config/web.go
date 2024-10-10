package config

import (
	"fmt"
	"time"
)

type web struct {
	bindingIpAddress       string
	port                   uint
	readTimeOutInSec       uint
	idleTimeoutInSec       uint
	writeTimeoutInSec      uint
	readHeaderTimeoutInSec uint
}

func (w web) Address() string {
	return fmt.Sprintf("%s:%d", w.bindingIpAddress, w.port)
}

func (w web) ReadTimeOut() time.Duration {
	return time.Duration(w.readTimeOutInSec) * time.Second
}

func (w web) IdleTimeout() time.Duration {
	return time.Duration(w.idleTimeoutInSec) * time.Second
}

func (w web) WriteTimeout() time.Duration {
	return time.Duration(w.writeTimeoutInSec) * time.Second
}

func (w web) ReadHeaderTimeout() time.Duration {
	return time.Duration(w.readHeaderTimeoutInSec) * time.Second
}
