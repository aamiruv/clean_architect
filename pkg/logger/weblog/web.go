package weblog

import (
	"bytes"
	"net/http"
)

type logger struct {
	url string
}

func New(url string) logger {
	return logger{url: url}
}

func (l logger) Write(p []byte) (int, error) {
	// we should ignore error to avoid panic
	/*res, err := http.Post(l.url, "application/json", bytes.NewBuffer(p))
	return res.StatusCode, err*/
	go http.Post(l.url, "application/json", bytes.NewBuffer(p))
	return 0, nil
}
