package logging

import (
	"os"
	"sync"

	"github.com/MadBase/MadNet/config"
	"github.com/MadBase/MadNet/constants"
	"github.com/sirupsen/logrus"
)

var loggers loggerDetails // singleton containing all loggers
type loggerDetails struct {
	sync.Once
	loggers map[string]*logrus.Logger
}

func Init(stdoutLogger *logrus.Logger) {
	loggers.Do(func() {

		f, err := openLogFile(config.Configuration.Logfile)
		if stdoutLogger != nil && err != nil {
			stdoutLogger.Errorf(err.Error())
		}

		minFileLvl, _ := logrus.ParseLevel(config.Configuration.Logfile.MinLevel)

		levels := config.LogLevelMap()
		loggers.loggers = make(map[string]*logrus.Logger, len(constants.ValidLoggers))
		for _, name := range constants.ValidLoggers {
			lvl, ok := levels[name]
			if !ok {
				lvl = logrus.InfoLevel
			}

			if stdoutLogger != nil {
				stdoutLogger.Infof("Setting log level for '%v' to '%v'...", name, lvl)
			}

			loggers.loggers[name] = MakeLogger(name, lvl, f, minFileLvl)
		}

		loggers.loggers["main"].Infof("Logging initialized successfully")
	})
}

func MakeStdOutLogger(loggerName string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&StdOutFormatter{Name: loggerName})
	return logger
}

func MakeLogger(loggerName string, lvl logrus.Level, file *os.File, minFileLvl logrus.Level) *logrus.Logger {
	stdOutLogger := MakeStdOutLogger(loggerName)
	stdOutLogger.SetLevel(lvl)
	if file == nil {
		return stdOutLogger
	}

	fileLogger := logrus.New()

	fileLoggerLvl := lvl
	if minFileLvl > lvl {
		fileLoggerLvl = minFileLvl
	}

	fileLogger.SetLevel(fileLoggerLvl)
	fileLogger.SetFormatter(&FileOutFormatter{Name: loggerName})
	fileLogger.SetOutput(file)
	fileLogger.AddHook(proxyLogger{stdOutLogger})

	return fileLogger
}
