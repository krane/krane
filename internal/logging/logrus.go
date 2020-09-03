package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// ConfigureLogrus : config a logrus logger
func ConfigureLogrus() {
	hostname, _ := os.Hostname()

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	// TODO: add contextual logger, these dont currently work since they arent bounded to a logger instance
	logrus.WithField("pid", os.Getpid())
	logrus.WithField("ppid", os.Getppid())
	logrus.WithField("hostname", hostname)
}
