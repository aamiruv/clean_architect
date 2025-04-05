package logger_test

import (
	"testing"

	"github.com/amirzayi/clean_architect/pkg/logger"
)

func TestWriteLogHourly(t *testing.T) {
	log := logger.NewFileLogger(logger.FileLogHourly, "./")
	if _, err := log.Write([]byte("test log hourly")); err != nil {
		t.Fail()
	}
}

func TestWriteLogDaily(t *testing.T) {
	log := logger.NewFileLogger(logger.FileLogDaily, "./")
	if _, err := log.Write([]byte("test log daily")); err != nil {
		t.Fail()
	}
}

func TestWriteLogMono(t *testing.T) {
	log := logger.NewFileLogger(logger.FileLogMono, "./")
	if _, err := log.Write([]byte("test log mono")); err != nil {
		t.Fail()
	}
}
