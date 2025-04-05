// Pacakge remotelog provides logged messages into files.
package remotelog

import (
	"bytes"
	"io"
	"net/http"
)

type logger struct {
	url string
}

func New(url string) io.Writer {
	return logger{url: url}
}

func (l logger) Write(p []byte) (int, error) {
	// we should ignore error to avoid panic
	/*res, err := http.Post(l.url, "application/json", bytes.NewBuffer(p))
	return res.StatusCode, err*/
	go http.Post(l.url, "application/json", bytes.NewBuffer(p))
	return 0, nil
}
