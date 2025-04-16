package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	logFile     *os.File
)

// LogParams is a map for structured logging
type LogParams map[string]interface{}

// InitializeLogger sets up the logger with the specified log directory
func InitializeLogger(logDir string) {
	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Create a log file with the current date
	currentTime := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("receipt_processor_%s.log", currentTime))

	var err error
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Create loggers with different prefixes for different log levels
	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	warnLogger = log.New(logFile, "WARN: ", log.Ldate|log.Ltime)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)

	LogInfo("Logger initialized", LogParams{"logPath": logPath})
}

// addFileInfo adds the file and line number to the log entry
func addFileInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	file = filepath.Base(file)
	return fmt.Sprintf("[%s:%d] ", file, line)
}

// formatParams formats the log parameters into a string
func formatParams(params LogParams) string {
	if len(params) == 0 {
		return ""
	}

	result := " {"
	for k, v := range params {
		result += fmt.Sprintf(" %s: %v,", k, v)
	}
	// Remove the trailing comma and add closing brace
	result = result[:len(result)-1] + " }"
	return result
}

// LogInfo logs an info message
func LogInfo(message string, params LogParams) {
	fileInfo := addFileInfo()
	formatted := formatParams(params)
	infoLogger.Println(fileInfo + message + formatted)

	fmt.Println("INFO: " + fileInfo + message + formatted)
}

// LogWarn logs a warning message
func LogWarn(message string, params LogParams) {
	fileInfo := addFileInfo()
	formatted := formatParams(params)
	warnLogger.Println(fileInfo + message + formatted)

	fmt.Println("WARN: " + fileInfo + message + formatted)
}

// LogError logs an error message
func LogError(message string, params LogParams) {
	fileInfo := addFileInfo()
	formatted := formatParams(params)
	errorLogger.Println(fileInfo + message + formatted)

	fmt.Println("ERROR: " + fileInfo + message + formatted)
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}
