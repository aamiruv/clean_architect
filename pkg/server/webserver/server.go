// Pacakge webserver provides http server construction.
package webserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func New(opts ...optionServerFunc) server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler:           mux,
		Addr:              address,
		MaxHeaderBytes:    maxHeaderBytes,
		IdleTimeout:       idleTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	httpServer := server{
		httpServer:      srv,
		shutdownTimeout: shutdownTimeout,
	}
	for _, opt := range opts {
		opt(&httpServer)
	}
	return httpServer
}

type optionServerFunc func(*server)

// WithHandler set handler to server
func WithHandler(handler http.Handler) optionServerFunc {
	return func(s *server) {
		s.httpServer.Handler = handler
	}
}

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

// WithTimeouts set idle, read, write, readHeader timeout values to server
func WithTimeouts(idle, read, write, readHeader, shutdown time.Duration) optionServerFunc {
	return func(s *server) {
		WithIdleTimeout(idle)(s)
		WithReadTimeout(read)(s)
		WithWriteTimeout(write)(s)
		WithReadHeaderTimeout(readHeader)(s)
		WithShutdownTimeout(shutdown)(s)
	}
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		err2 := s.httpServer.Close()
		if err2 != nil {
			return errors.Join(err, err2)
		}
		return err
	}
	return nil
}
