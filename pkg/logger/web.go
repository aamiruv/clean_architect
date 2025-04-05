// Package logger provides logged messages into files and remote http logger.
package logger

import (
	"bytes"
	"io"
	"net/http"
)

type remoteLogger struct {
	url string
}

func NewRemoteLogger(url string) io.Writer {
	return &remoteLogger{url: url}
}

func (l *remoteLogger) Write(p []byte) (int, error) {
	// we should ignore error to avoid panic

	go http.Post(l.url, "application/json", bytes.NewBuffer(p))
	return 0, nil
}
