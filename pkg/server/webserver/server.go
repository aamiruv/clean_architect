package webserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

func New(opts ...optionServerFunc) server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler:           mux,
		Addr:              address,
		ErrorLog:          logger,
		MaxHeaderBytes:    maxHeaderBytes,
		IdleTimeout:       idleTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	httpServer := server{httpServer: srv, mux: mux}
	for _, opt := range opts {
		opt(&httpServer)
	}
	return httpServer
}

type optionServerFunc func(*server)

func WithAddress(address string) optionServerFunc {
	return func(s *server) {
		s.httpServer.Addr = address
	}
}

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

func WithTimeout(idle, read, write, readHeader time.Duration) optionServerFunc {
	return func(s *server) {
		s.httpServer.IdleTimeout = idle
		s.httpServer.ReadTimeout = read
		s.httpServer.WriteTimeout = write
		s.httpServer.ReadHeaderTimeout = readHeader
	}
}

func (s *server) MuxHandler() *http.ServeMux {
	return s.mux
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) GracefulShutdown(deadline time.Duration) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(deadline))
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		err2 := s.httpServer.Close()
		if err2 != nil {
			return errors.Join(err, err2)
		}
	}
	return nil
}
