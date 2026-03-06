package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
func formatAndSanitize(v ...interface{}) string {
	var builder strings.Builder
	for i, arg := range v {
		if i > 0 {
			builder.WriteString(" ")
		}
		if s, ok := arg.(string); ok {
			builder.WriteString(strconv.Quote(s))
		} else {
			builder.WriteString(fmt.Sprint(arg))
		}
	}
	return builder.String()
}
func LogDebug(v ...interface{}) {
	if currentLogLevel >= LevelDebug {
		log.Println(formatAndSanitize(v...))
	}
}
func LogInfo(v ...interface{}) {
	if currentLogLevel >= LevelInfo {
		log.Println(formatAndSanitize(v...))
	}
}
func LogWarn(v ...interface{}) {
	if currentLogLevel >= LevelWarn {
		log.Println(formatAndSanitize(v...))
	}
}
func LogError(v ...interface{}) {
	if currentLogLevel >= LevelError {
		log.Println(formatAndSanitize(v...))
	}
}
