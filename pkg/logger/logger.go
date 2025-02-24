// pkg/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"log"
	//"os"
)

// Level represents the severity of a log message.
type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger is a structured logger with support for log levels.
type Logger struct {
	level  Level
	logger *log.Logger
}

// New creates a new Logger with the specified log level and output.
func New(level Level, out io.Writer) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(out, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Debugf logs a debug message.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= Debug {
		l.logf("DEBUG", format, v...)
	}
}

// Infof logs an info message.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= Info {
		l.logf("INFO", format, v...)
	}
}

// Warnf logs a warning message.
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= Warn {
		l.logf("WARN", format, v...)
	}
}

// Errorf logs an error message.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= Error {
		l.logf("ERROR", format, v...)
	}
}

// Panicf logs a message and panics.
func (l *Logger) Panicf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Output(2, fmt.Sprintf("PANIC: %s", msg))
	panic(msg)
}

// logf formats the log message with the level prefix.
func (l *Logger) logf(level string, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.logger.Output(3, fmt.Sprintf("[%s] %s", level, msg))
}

// ParseLevel converts a string to a log Level (case-insensitive).
func ParseLevel(s string) (Level, error) {
	switch s {
	case "debug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error":
		return Error, nil
	default:
		return Info, fmt.Errorf("invalid log level: %s", s)
	}
}