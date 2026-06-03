package logger

import (
	"log"
	"os"
)

var (
	debugMode bool
	infoLog   = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	errorLog  = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	warnLog   = log.New(os.Stderr, "[WARN] ", log.LstdFlags)
)

// SetDebug enables or disables debug mode.
// When disabled, Info messages are suppressed.
func SetDebug(enabled bool) {
	debugMode = enabled
}

// IsDebug returns whether debug mode is active.
func IsDebug() bool {
	return debugMode
}

// Info logs an informational message (only when debug mode is enabled).
func Info(format string, v ...any) {
	if debugMode {
		infoLog.Printf(format, v...)
	}
}

// Warn logs a warning message (always shown).
func Warn(format string, v ...any) {
	warnLog.Printf(format, v...)
}

// Error logs an error message (always shown).
func Error(format string, v ...any) {
	errorLog.Printf(format, v...)
}
