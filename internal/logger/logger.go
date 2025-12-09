package logger

import (
	"log/slog"
	"os"
	"strings"
)

func Initlogger(logLevel string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(logLevel),
	})
	// create the logger
	logger := slog.New(handler)

	return logger
}

func getLogLevel(logLevel string) slog.Level {
	var levelInfo slog.Level
	switch strings.ToLower(logLevel) {
	case "debug": //priority is low (if someone pass this then it will print all debug,info,warn,error)
		levelInfo = slog.LevelDebug
	case "info": //(if someone pass this then it will print info, warn, error)
		levelInfo = slog.LevelInfo
	case "warn": //(if someone pass this then it will print warn, error)
		levelInfo = slog.LevelWarn
	case "error": //priority is high (if someone pass this then only it will print error log only)
		levelInfo = slog.LevelError
	default:
		levelInfo = slog.LevelInfo
	}
	return levelInfo
}
