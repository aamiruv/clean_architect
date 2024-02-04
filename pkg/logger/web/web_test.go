package web

import "testing"

func TestWriteLogOnWeb(t *testing.T) {
	logger := NewLogger("https://google.com")
	if _, err := logger.Write([]byte("hi google!")); err != nil {
		t.Fail()
	}
}
