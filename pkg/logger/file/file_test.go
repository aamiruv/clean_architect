package file

import "testing"

func TestWriteLogHourly(t *testing.T) {
	logger := NewLogger(LogHourly, "./")
	if _, err := logger.Write([]byte("test log hourly")); err != nil {
		t.Fail()
	}
}

func TestWriteLogDaily(t *testing.T) {
	logger := NewLogger(LogDaily, "./")
	if _, err := logger.Write([]byte("test log daily")); err != nil {
		t.Fail()
	}
}

func TestWriteLogMono(t *testing.T) {
	logger := NewLogger(LogMono, "./")
	if _, err := logger.Write([]byte("test log mono")); err != nil {
		t.Fail()
	}
}
