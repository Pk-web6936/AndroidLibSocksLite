package logging

import (
	"github.com/sirupsen/logrus"
	"os"
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
		"file":  "example.go", // Example: Adding file name
	}).Log(level, message)
}

func init() {
	// Set log format to JSON
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// Set default log level
	logrus.SetLevel(logrus.InfoLevel)
	
	// File logging hook
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600) // <-- Changed from 0666 to 0600
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(file)
}
