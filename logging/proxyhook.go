package logging

import "github.com/sirupsen/logrus"

type proxyLogger struct {
	proxyLogger *logrus.Logger
}

func (proxyLogger) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (p proxyLogger) Fire(e *logrus.Entry) error {
	p.proxyLogger.
		WithFields(e.Data).
		WithTime(e.Time).
		Log(e.Level, e.Message)

	return nil
}
