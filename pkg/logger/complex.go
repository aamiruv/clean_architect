package logger

import (
	"io"
)

type ComplexLogger struct {
	logger []io.Writer
}

func NewComplexLogger(loggers ...io.Writer) ComplexLogger {
	return ComplexLogger{logger: loggers}
}

func (l ComplexLogger) Write(p []byte) (int, error) {
	for _, logger := range l.logger {
		go logger.Write(p)
	}
	return 0, nil
}
