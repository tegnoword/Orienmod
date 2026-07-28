package logger

import (
	"log"
	"os"
)

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type Logger struct {
	level Level
}

func NewLogger() *Logger {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	return &Logger{level: Level(level)}
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	if l.level == LevelDebug || l.level == LevelInfo {
		log.Printf("[INFO] %s %v", msg, fields)
	}
}

func (l *Logger) Error(msg string, err error) {
	log.Printf("[ERROR] %s: %v", msg, err)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	if l.level == LevelDebug {
		log.Printf("[DEBUG] %s %v", msg, fields)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level == LevelDebug || l.level == LevelInfo || l.level == LevelWarn {
		log.Printf("[WARN] %s", msg)
	}
}

func (l *Logger) Request(method, path, email string) {
	l.Debug("📨 Petición", "method", method, "path", path, "email", email)
}
