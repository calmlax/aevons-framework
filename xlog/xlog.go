// Package xlog provides a lightweight wrapper around the standard log package.
// Replace the implementation here with zap/zerolog/slog as needed.
package xlog

import (
	"fmt"
	"log"
)

// Info logs an informational message.
func Info(msg string, args ...any) {
	log.Printf("[INFO]  "+msg, args...)
}

// Warn logs a warning message.
func Warn(msg string, args ...any) {
	log.Printf("[WARN]  "+msg, args...)
}

// Error logs an error message.
func Error(msg string, args ...any) {
	log.Printf("[ERROR] "+msg, args...)
}

// Fatal logs a fatal message and exits.
func Fatal(msg string, args ...any) {
	log.Fatalf("[FATAL] "+msg, args...)
}

// Debug logs a debug message.
func Debug(msg string, args ...any) {
	log.Printf("[DEBUG] "+msg, args...)
}

// Infof is an alias for Info with fmt-style formatting.
func Infof(format string, args ...any) { Info(fmt.Sprintf(format, args...)) }

// Errorf is an alias for Error with fmt-style formatting.
func Errorf(format string, args ...any) { Error(fmt.Sprintf(format, args...)) }
