package mocks

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

func NewMockLogger() *logrus.Logger {
	logger := logrus.New()
	if os.Getenv("TEST_DEBUG") == "" {
		logger.SetOutput(ioutil.Discard)
	}
	return logger
}
