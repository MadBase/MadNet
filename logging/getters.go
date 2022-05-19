package logging

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

//GetLogger either returns an existing logger for package specified or creates a new one
func GetLogger(name string) *logrus.Logger {
	Init(nil)

	logger, exists := loggers.loggers[strings.ToLower(name)]
	if !exists {
		panic(fmt.Sprintf("Invalid logger requested: %v", name))
	}

	return logger
}

//GetLogWriter returns an io.Writer that maps to the named logger at the specified level
func GetLogWriter(pkgName string, level logrus.Level) *LogWriter {
	Init(nil)

	return &LogWriter{GetLogger(pkgName), level}
}

//LogWriter struct used to provide an io.Writer
type LogWriter struct {
	logger *logrus.Logger
	level  logrus.Level
}

func (logWriter *LogWriter) Write(p []byte) (n int, err error) {
	logWriter.logger.Log(logWriter.level, strings.TrimRight(string(p), "\n"))
	return len(p), nil
}
