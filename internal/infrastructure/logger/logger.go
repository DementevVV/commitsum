// Package logger provides logging functionality.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Level represents a log level.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a configurable logger.
type Logger struct {
	logger   *log.Logger
	level    Level
	file     *os.File
	disabled bool
}

var defaultLogger *Logger

// Init initializes the logger.
func Init(level Level, version, buildTime string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	logDir := filepath.Join(homeDir, ".config", "commitsum", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFileName := fmt.Sprintf("commitsum-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	writers := []io.Writer{file}
	if level == LevelDebug {
		writers = append(writers, os.Stderr)
	}

	multiWriter := io.MultiWriter(writers...)
	logger := log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	defaultLogger = &Logger{
		logger:   logger,
		level:    level,
		file:     file,
		disabled: false,
	}

	defaultLogger.Info("CommitSum started", "version", version, "build_time", buildTime)

	return nil
}

// Disable disables logging (for tests).
func Disable() {
	if defaultLogger != nil {
		defaultLogger.disabled = true
	}
}

// Enable enables logging.
func Enable() {
	if defaultLogger != nil {
		defaultLogger.disabled = false
	}
}

// Close closes the logger and file.
func Close() error {
	if defaultLogger != nil && defaultLogger.file != nil {
		defaultLogger.Info("CommitSum shutting down")
		return defaultLogger.file.Close()
	}
	return nil
}

// logMessage writes a message with the specified level.
func (l *Logger) logMessage(level Level, msg string, keyvals ...interface{}) {
	if l.disabled || level < l.level {
		return
	}

	var kvStr string
	if len(keyvals) > 0 {
		kvStr = " |"
		for i := 0; i < len(keyvals); i += 2 {
			if i+1 < len(keyvals) {
				kvStr += fmt.Sprintf(" %v=%v", keyvals[i], keyvals[i+1])
			} else {
				kvStr += fmt.Sprintf(" %v=%v", keyvals[i], "")
			}
		}
	}

	l.logger.Printf("[%s] %s%s", level.String(), msg, kvStr)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string, keyvals ...interface{}) {
	l.logMessage(LevelDebug, msg, keyvals...)
}

// Info logs an info message.
func (l *Logger) Info(msg string, keyvals ...interface{}) {
	l.logMessage(LevelInfo, msg, keyvals...)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string, keyvals ...interface{}) {
	l.logMessage(LevelWarn, msg, keyvals...)
}

// Error logs an error message.
func (l *Logger) Error(msg string, keyvals ...interface{}) {
	l.logMessage(LevelError, msg, keyvals...)
}

// Global logging functions.

// Debug logs a debug message.
func Debug(msg string, keyvals ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, keyvals...)
	}
}

// Info logs an info message.
func Info(msg string, keyvals ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, keyvals...)
	}
}

// Warn logs a warning message.
func Warn(msg string, keyvals ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, keyvals...)
	}
}

// Error logs an error message.
func Error(msg string, keyvals ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, keyvals...)
	}
}

// LogGitHubCommand logs a GitHub CLI command execution.
func LogGitHubCommand(cmd string, duration time.Duration, err error) {
	if err != nil {
		Error("GitHub CLI command failed",
			"command", cmd,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error())
	} else {
		Info("GitHub CLI command executed",
			"command", cmd,
			"duration_ms", duration.Milliseconds())
	}
}

// LogUserAction logs a user action.
func LogUserAction(action string, details ...interface{}) {
	Info("User action: "+action, details...)
}

// LogPerformance logs performance metrics.
func LogPerformance(operation string, duration time.Duration, details ...interface{}) {
	allDetails := append([]interface{}{"duration_ms", duration.Milliseconds()}, details...)
	Info("Performance: "+operation, allDetails...)
}
