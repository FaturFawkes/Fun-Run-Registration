package utils

import (
	"fmt"
	"log"
	"time"
)

// Logger provides structured logging
type Logger struct {
	prefix string
}

// NewLogger creates a new logger with a prefix
func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Info logs informational messages
func (l *Logger) Info(message string, args ...interface{}) {
	l.log("INFO", message, args...)
}

// Error logs error messages
func (l *Logger) Error(message string, args ...interface{}) {
	l.log("ERROR", message, args...)
}

// Warning logs warning messages
func (l *Logger) Warning(message string, args ...interface{}) {
	l.log("WARN", message, args...)
}

// Debug logs debug messages
func (l *Logger) Debug(message string, args ...interface{}) {
	l.log("DEBUG", message, args...)
}

// log formats and outputs log message
func (l *Logger) log(level, message string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	prefix := ""
	if l.prefix != "" {
		prefix = fmt.Sprintf("[%s] ", l.prefix)
	}
	
	formattedMessage := fmt.Sprintf(message, args...)
	log.Printf("[%s] %s%s: %s", timestamp, prefix, level, formattedMessage)
}

// Global logger instances for convenience
var (
	ServerLogger = NewLogger("SERVER")
	DBLogger     = NewLogger("DATABASE")
	AuthLogger   = NewLogger("AUTH")
	EmailLogger  = NewLogger("EMAIL")
)
