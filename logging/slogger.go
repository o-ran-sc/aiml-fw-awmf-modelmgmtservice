package logging

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"

	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logLevel string

var (
	LOG_LEVEL_DEBUG logLevel = "DEBUG"
	LOG_LEVEL_INFO  logLevel = "INFO"
	LOG_LEVEL_WARN  logLevel = "WARN"
	LOG_LEVEL_ERROR logLevel = "ERROR"
)

func (l logLevel) convertLogLevelToSlogLogLevel() slog.Level {
	switch l {
	case LOG_LEVEL_DEBUG:
		return slog.LevelDebug
	case LOG_LEVEL_INFO:
		return slog.LevelInfo
	case LOG_LEVEL_WARN:
		return slog.LevelWarn
	case LOG_LEVEL_ERROR:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func parseLogLevel(level string) (logLevel, error) {
	switch strings.ToUpper(level) {
	case string(LOG_LEVEL_DEBUG), string(LOG_LEVEL_INFO), string(LOG_LEVEL_WARN), string(LOG_LEVEL_ERROR):
		return logLevel(level), nil
	default:
		return LOG_LEVEL_INFO, errors.New("invalid log level, set default log level 'INFO'")
	}
}

var (
	Logger *slog.Logger
	once   sync.Once
)

func Load(logLevel string, filename string) {
	once.Do(func() {
		parsedLogLevel, err := parseLogLevel(logLevel)
		fileRotationLogger := lumberjack.Logger{
			Filename:   filename,
			MaxSize:    100, // 100MB
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		}

		logHandlerOpt := slog.HandlerOptions{
			Level:     parsedLogLevel.convertLogLevelToSlogLogLevel(),
			AddSource: true,
		}

		fileLogHandler := slog.NewJSONHandler(&fileRotationLogger, &logHandlerOpt)
		stdoutLogHandler := slog.NewTextHandler(os.Stdout, &logHandlerOpt)

		Logger = slog.New(slogmulti.Fanout(fileLogHandler, stdoutLogHandler)).With(slog.String("app", "mmes"))
		if err != nil {
			Logger.Error("error occurred: ", slog.Any("error", err))
		}
	})
}
