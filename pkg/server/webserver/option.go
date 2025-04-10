package webserver

import (
	"log"
	"time"
)

type optionServerFunc func(*server)

// WithAddress set access address to server
func WithAddress(address string) optionServerFunc {
	return func(s *server) {
		s.httpServer.Addr = address
	}
}

// WithLogger set log.Logger to the server error logger
func WithLogger(logger *log.Logger) optionServerFunc {
	return func(s *server) {
		s.httpServer.ErrorLog = logger
	}
}

func WithMaxHeaderBytes(bytes int) optionServerFunc {
	return func(s *server) {
		s.httpServer.MaxHeaderBytes = bytes
	}
}

func WithIdleTimeout(idle time.Duration) optionServerFunc {
	return func(s *server) {
		s.httpServer.IdleTimeout = idle
	}
}

func WithReadTimeout(read time.Duration) optionServerFunc {
	return func(s *server) {
		s.httpServer.ReadTimeout = read
	}
}

func WithWriteTimeout(write time.Duration) optionServerFunc {
	return func(s *server) {
		s.httpServer.ReadTimeout = write
	}
}

func WithReadHeaderTimeout(read time.Duration) optionServerFunc {
	return func(s *server) {
		s.httpServer.ReadTimeout = read
	}
}

func WithShutdownTimeout(shutdown time.Duration) optionServerFunc {
	return func(s *server) {
		s.shutdownTimeout = shutdown
	}
}

// WithTimeouts set idle, read, write, readHeader, shutdown timeout values to server
func WithTimeouts(idle, read, write, readHeader, shutdown time.Duration) optionServerFunc {
	return func(s *server) {
		WithIdleTimeout(idle)(s)
		WithReadTimeout(read)(s)
		WithWriteTimeout(write)(s)
		WithReadHeaderTimeout(readHeader)(s)
		WithShutdownTimeout(shutdown)(s)
	}
}
