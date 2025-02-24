package Common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// LogEntry represents a structured log entry.
type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// LogInfo logs an info message.
func LogInfo(message string, fields map[string]interface{}) {
	logMessage("INFO", message, fields)
}

// LogWarning logs a warning message.
func LogWarning(message string, fields map[string]interface{}) {
	logMessage("WARNING", message, fields)
}

// LogError logs an error message.
func LogError(message string, err error, fields map[string]interface{}) {
	if err != nil {
		if fields == nil {
			fields = make(map[string]interface{})
		}
		fields["error"] = err.Error()
	}
	logMessage("ERROR", message, fields)
}

// logMessage is a helper function to log messages in a structured way.
func logMessage(level, message string, fields map[string]interface{}) {
	entry := LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    fields,
	}

	// Convert log entry to JSON for structured logging
	logData, err := json.Marshal(entry)
	if err != nil {
		// If marshaling to JSON fails, fallback to simple log.
		log.Printf("[%s] %s - Error marshaling log entry: %v", level, message, err)
		return
	}

	// Write the log entry to the log output.
	log.Printf(string(logData))
}

// LogDuration logs the duration of an operation.
func LogDuration(operation string, startTime time.Time, fields map[string]interface{}) {
	duration := time.Since(startTime).Seconds()
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["duration"] = duration
	LogInfo(fmt.Sprintf("%s completed", operation), fields)
}
