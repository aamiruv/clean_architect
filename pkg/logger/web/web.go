package web

import (
	"bytes"
	"net/http"
)

type Logger struct {
	URL string
}

func NewLogger(url string) Logger {
	return Logger{URL: url}
}

func (l Logger) Write(p []byte) (int, error) {
	// we should ignore error to avoid panic
	/*res, err := http.Post(l.URL, "application/json", bytes.NewBuffer(p))
	return res.StatusCode, err*/
	go http.Post(l.URL, "application/json", bytes.NewBuffer(p))
	return 0, nil
}
