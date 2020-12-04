package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/constants"
)

var once sync.Once
var l *logrus.Logger

func level() logrus.Level {
	level, err := logrus.ParseLevel(os.Getenv(constants.EnvLogLevel))
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

func Error(err error) {
	l.Log(logrus.ErrorLevel, err)
}

func Errorf(format string, err error) {
	l.Logf(logrus.ErrorLevel, format, err)
}

func Debug(msg string) {
	l.Log(logrus.DebugLevel, msg)
}

func Debugf(format string, args ...interface{}) {
	l.Logf(logrus.DebugLevel, format, args)
}

func Info(msg string) {
	l.Log(logrus.InfoLevel, msg)
}

func Infof(format string, args ...interface{}) {
	l.Logf(logrus.InfoLevel, format, args)
}

func Warn(msg string) {
	l.Log(logrus.WarnLevel, msg)
}

func Warnf(format string, args ...interface{}) {
	l.Logf(logrus.WarnLevel, format, args)
}

func Trace(msg string) {
	l.Log(logrus.TraceLevel, msg)
}

func Tracef(format string, args ...interface{}) {
	l.Logf(logrus.TraceLevel, format, args)
}

func Fatal(msg string) {
	l.Log(logrus.FatalLevel, msg)
}

func Fatalf(format string, args ...interface{}) {
	l.Logf(logrus.FatalLevel, format, args)
}
