// Pacakge filelog provides store logged messages into file(s).
package filelog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LoggerType uint8

const (
	logNone LoggerType = iota
	LogMono
	LogHourly
	LogDaily
)

type Logger struct {
	loggerType LoggerType
	directory  string
}

func New(loggerType LoggerType, directory string) Logger {
	return Logger{
		loggerType: loggerType,
		directory:  directory,
	}
}

func (l Logger) Write(p []byte) (int, error) {
	var logFileName string

	y, m, d := time.Now().Date()

	switch l.loggerType {
	case logNone:
		return 0, nil
	case LogHourly:
		h := time.Now().Hour()
		logFileName = fmt.Sprintf("%s/%d/%d/%d/%d.log", l.directory, y, m, d, h)
	case LogDaily:
		logFileName = fmt.Sprintf("%s/%d/%d/%d.log", l.directory, y, m, d)
	case LogMono:
		logFileName = fmt.Sprintf("%s/log.log", l.directory)

	default:
		return 0, fmt.Errorf("invalid logger type")
	}

	file, err := openLogFile(logFileName)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = file.Close()
	}()
	return file.Write(p)
}

func openLogFile(filePath string) (*os.File, error) {
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
