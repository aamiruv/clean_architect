package filelog

import "testing"

func TestWriteLogHourly(t *testing.T) {
	logger := New(LogHourly, "./")
	if _, err := logger.Write([]byte("test log hourly")); err != nil {
		t.Fail()
	}
}

func TestWriteLogDaily(t *testing.T) {
	logger := New(LogDaily, "./")
	if _, err := logger.Write([]byte("test log daily")); err != nil {
		t.Fail()
	}
}

func TestWriteLogMono(t *testing.T) {
	logger := New(LogMono, "./")
	if _, err := logger.Write([]byte("test log mono")); err != nil {
		t.Fail()
	}
}
