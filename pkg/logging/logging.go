package logging

import (
	"github.com/sirupsen/logrus"
)

// LogInfo logs informational messages with an [INFO] prefix.
func LogInfo(message string) {
	logMessage(logrus.InfoLevel, message)
}

// LogError logs error messages with an [ERROR] prefix.
func LogError(message string) {
	logMessage(logrus.ErrorLevel, message)
}

// logMessage is a helper function to format and log messages.
func logMessage(level logrus.Level, message string) {
	logrus.WithFields(logrus.Fields{
		"level": level.String(),
	}).Log(level, message)
}

func init() {
	// Set log format to JSON for better structure and parsing
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Set log level to info as default
	logrus.SetLevel(logrus.InfoLevel)
}
