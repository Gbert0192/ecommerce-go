package log

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func SetupLoger() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true, // ada log info error and warning
		FullTimestamp: true,
	})

	log.Info("Logged initiated using logrus!")
	Logger = log
}
