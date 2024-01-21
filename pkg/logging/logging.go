package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

// New make logger
func New() *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
