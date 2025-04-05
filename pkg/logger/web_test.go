package logger_test

import (
	"testing"

	"github.com/amirzayi/clean_architect/pkg/logger"
)

func TestWriteLogOnWeb(t *testing.T) {
	log := logger.NewRemoteLogger("https://google.com")
	if _, err := log.Write([]byte("hi google!")); err != nil {
		t.Fail()
	}
}
