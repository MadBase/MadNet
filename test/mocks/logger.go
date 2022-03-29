package mocks

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func NewMockLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)
	return logger
}
