package webserver

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	const (
		timeout                 = time.Second
		mockedHttpServerAddress = "localhost:54321"
		testingRoute            = "/test"
		fullTestingRoute        = "http://" + mockedHttpServerAddress + testingRoute
	)

	mux := http.NewServeMux()
	srv := New(WithAddress(mockedHttpServerAddress), WithHandler(mux))

	t.Run("test run http server", func(t *testing.T) {
		go func() {
			if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				t.Fail()
			}
		}()
		mux.HandleFunc(testingRoute, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	t.Run("test http server functionality", func(t *testing.T) {
		res, err := http.Get(fullTestingRoute)
		if err != nil || res.StatusCode != http.StatusOK {
			t.Fail()
		}
	})

	t.Run("test shut down gracefully http server", func(t *testing.T) {
		if err := srv.GracefulShutdown(); err != nil {
			t.Fail()
		}
	})
}
