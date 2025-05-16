package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipCompress(level int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}
			gzw, err := gzip.NewWriterLevel(w, level)
			if err != nil {
				panic(err)
			}
			w.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(gzipWriter{w, gzw}, r)
			gzw.Close()
		})
	}
}

type gzipWriter struct {
	http.ResponseWriter
	w io.Writer
}

func (g gzipWriter) Write(p []byte) (int, error) {
	return g.w.Write(p)
}
