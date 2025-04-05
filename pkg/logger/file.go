// Pacakge filelog provides store logged messages into file(s).
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileLoggerType uint8

const (
	logNone FileLoggerType = iota
	FileLogMono
	FileLogHourly
	FileLogDaily
)

type fileLogger struct {
	loggerType FileLoggerType
	directory  string
}

func NewFileLogger(loggerType FileLoggerType, directory string) io.Writer {
	return &fileLogger{
		loggerType: loggerType,
		directory:  directory,
	}
}

func (l *fileLogger) Write(p []byte) (int, error) {
	var logFileName string

	y, m, d := time.Now().Date()

	switch l.loggerType {
	case logNone:
		return 0, nil
	case FileLogHourly:
		h := time.Now().Hour()
		logFileName = fmt.Sprintf("%s/%d/%d/%d/%d.log", l.directory, y, m, d, h)
	case FileLogDaily:
		logFileName = fmt.Sprintf("%s/%d/%d/%d.log", l.directory, y, m, d)
	case FileLogMono:
		logFileName = fmt.Sprintf("%s/log.log", l.directory)

	default:
		return 0, fmt.Errorf("invalid logger type")
	}

	file, err := l.openLogFile(logFileName)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()
	return file.Write(p)
}

func (l *fileLogger) openLogFile(filePath string) (*os.File, error) {
	if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}
