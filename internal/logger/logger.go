package logger

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/constants"
)

// ERROR : log in error format
func ERROR(msg string) {
}

// Configure : configure the logrus logger
func Configure() {
	logLevel, err := logrus.ParseLevel(os.Getenv(constants.EnvLogLevel))
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: false})

	// TODO: add contextual logger, these dont currently work since they arent bounded to a logger instance
	hostname, _ := os.Hostname()
	logrus.WithField("pid", os.Getpid())
	logrus.WithField("ppid", os.Getppid())
	logrus.WithField("hostname", hostname)
}
