package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// ConfigureLogrus : config a logrus logger
func ConfigureLogrus() {
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: false})
}
