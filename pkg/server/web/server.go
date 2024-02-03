package web

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

func NewServer(
	address string,
	logger *log.Logger,
	maxHeaderBytes int,
	idleTimeout,
	readTimeout,
	writeTimeout,
	readHeaderTimeout time.Duration,
) server {
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

	return server{httpServer: srv, mux: mux}
}

func (s *server) GetMuxHandler() *http.ServeMux {
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
