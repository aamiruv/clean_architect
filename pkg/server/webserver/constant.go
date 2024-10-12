package webserver

import (
	"time"
)

const (
	address           = "0.0.0.0:8071"
	maxHeaderBytes    = 1 << 7
	idleTimeout       = 10 * time.Second
	readTimeout       = 7 * time.Second
	writeTimeout      = 15 * time.Second
	readHeaderTimeout = 1 * time.Second
)
