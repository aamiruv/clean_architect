// Pacakge webserver provides http server construction.
package webserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func New(handler http.Handler, options ...optionServerFunc) server {
	srv := &http.Server{
		Handler:           handler,
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
	for _, option := range options {
		option(&httpServer)
	}
	return httpServer
}

func (s *server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *server) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		if err2 := s.httpServer.Close(); err2 != nil {
			return errors.Join(err, err2)
		}
		return err
	}
	return nil
}
