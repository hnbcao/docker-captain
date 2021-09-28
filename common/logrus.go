package common

import "github.com/sirupsen/logrus"

func SetLogLevel(level string) {
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}
}
