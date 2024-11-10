package logging

import "log"

// LogInfo provides informational logging
func LogInfo(message string) {
	log.Printf("[INFO] %s\n", message)
}

// LogError provides error logging
func LogError(message string) {
	log.Printf("[ERROR] %s\n", message)
}
