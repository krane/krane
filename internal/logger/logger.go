package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/utils"
)

var once sync.Once
var l *logrus.Logger

func level() logrus.Level {
	level, err := logrus.ParseLevel(utils.EnvOrDefault(constants.EnvLogLevel, "info"))
	if err != nil {
		panic(err)
	}
	return level
}

func Configure() {
	once.Do(func() {
		l = logrus.New()
		l.SetOutput(os.Stdout)
		l.SetLevel(level())
		l.SetFormatter(&logrus.JSONFormatter{PrettyPrint: false})
	})
}

// withContext : add context to the logs
func withContext() *logrus.Entry {
	// if logger has not been configured, configure it
	if l == nil {
		Configure()
	}

	hostname, _ := os.Hostname()
	return l.
		WithField("hostname", hostname).
		WithField("pid", os.Getpid())
}

func Error(err error) {
	withContext().Error(err)
}

func Errorf(format string, err error) {
	withContext().Errorf(format, err)
}

func Debug(msg string) {
	withContext().Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	withContext().Debugf(format, args)
}

func Info(msg string) {
	withContext().Info(msg)
}

func Infof(format string, args ...interface{}) {
	withContext().Infof(format, args)
}

func Warn(msg string) {
	withContext().Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	withContext().Warnf(format, args)
}

func Trace(msg string) {
	withContext().Trace(msg)
}

func Tracef(format string, args ...interface{}) {
	withContext().Trace(format, args)
}

func Fatal(msg string) {
	withContext().Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	withContext().Fatalf(format, args)
}
