package util

import "log/slog"

func Debug(logger *slog.Logger, msg string, args ...any) {
	if logger != nil {
		logger.Debug(msg, args...)
		return
	}

	slog.Debug(msg, args...)
}

func Info(logger *slog.Logger, msg string, args ...any) {
	if logger != nil {
		logger.Info(msg, args...)
		return
	}

	slog.Info(msg, args...)
}

func Warn(logger *slog.Logger, msg string, args ...any) {
	if logger != nil {
		logger.Warn(msg, args...)
		return
	}

	slog.Warn(msg, args...)
}

func Error(logger *slog.Logger, msg string, args ...any) {
	if logger != nil {
		logger.Error(msg, args...)
		return
	}

	slog.Error(msg, args...)
}
