package main

import (
	"log"
	"os"
)

var logLevel = "INFO" // Default to INFO level

func init() {
	// Set log flags to include timestamp and file location
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Optionally read from environment variable
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		logLevel = level
	}
}

func logDebug(format string, v ...interface{}) {
	if logLevel == "DEBUG" {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func logInfo(format string, v ...interface{}) {
	if logLevel == "DEBUG" || logLevel == "INFO" {
		log.Printf("[INFO] "+format, v...)
	}
}

func logError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}
