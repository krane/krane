package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var standardLogger *logrus.Logger

// NewLogger :
func NewLogger() *logrus.Logger {
	if standardLogger != nil {
		return standardLogger
	}

	standardLogger = logrus.New()

	// set log level
	switch os.Getenv("KRANE_LOG_LEVEL") {
	case "info":
		standardLogger.SetLevel(logrus.InfoLevel)
	case "debug":
		standardLogger.SetLevel(logrus.DebugLevel)
	default:
		standardLogger.SetLevel(logrus.InfoLevel)
	}

	return standardLogger
}

// Debug :
func Debug(message string) {
	standardLogger.Debug(message)
}

// Debugf :
func Debugf(message string, args ...interface{}) {
	standardLogger.Debugf(message, args...)
}

// Info :
func Info(message string) {
	standardLogger.Info(message)
}

// Infof :
func Infof(message string, args ...interface{}) {
	standardLogger.Infof(message, args...)
}

// Error :
func Error(message string) {
	standardLogger.Error(message)
}

// Errorf :
func Errorf(message string, args ...interface{}) {
	standardLogger.Errorf(message, args...)
}
