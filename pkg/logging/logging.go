package logging

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var instance *logrus.Logger
var once sync.Once

// GetLogger returns singleton logrus logger
func GetLogger() *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()
	})
	return instance
}
