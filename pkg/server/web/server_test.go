package web

import (
	"net/http"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	timeout := time.Second
	mockedHttpServerAddress := "localhost:54321"
	testingRoute := "/test"
	fullTestingRoute := "http://" + mockedHttpServerAddress + testingRoute

	srv := NewServer(mockedHttpServerAddress, nil, 100000, timeout, timeout, timeout, timeout)
	go func() {
		if err := srv.Run(); err != nil {
			t.Fail()
		}
	}()

	srv.GetMuxHandler().HandleFunc(testingRoute, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	res, err := http.Get(fullTestingRoute)
	if err != nil {
		t.Fail()
	}
	if res.StatusCode != http.StatusOK {
		t.Fail()
	}

	time.Sleep(timeout)
	if err := srv.GracefulShutdown(10 * time.Second); err != nil {
		t.Error(err)
	}

}
