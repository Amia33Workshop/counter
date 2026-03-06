package main

import (
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LevelNone LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var currentLogLevel = LevelInfo

func InitLogger() {
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch levelStr {
	case "none":
		currentLogLevel = LevelNone
	case "error":
		currentLogLevel = LevelError
	case "warn":
		currentLogLevel = LevelWarn
	case "info":
		currentLogLevel = LevelInfo
	case "debug":
		currentLogLevel = LevelDebug
	}
}
func LogDebug(v ...interface{}) {
	if currentLogLevel >= LevelDebug {
		log.Println(v...)
	}
}
func LogInfo(v ...interface{}) {
	if currentLogLevel >= LevelInfo {
		log.Println(v...)
	}
}
func LogWarn(v ...interface{}) {
	if currentLogLevel >= LevelWarn {
		log.Println(v...)
	}
}
func LogError(v ...interface{}) {
	if currentLogLevel >= LevelError {
		log.Println(v...)
	}
}
