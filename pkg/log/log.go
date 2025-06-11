package log

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
	"time"
)

const (
	duration = "duration"
	empty    = "empty"
	errStr   = "error"

	traceLevel string = "trace"
	debugLevel string = "debug"
	infoLevel  string = "info"
	warnLevel  string = "warn"
	errorLevel string = "error"

	TraceLevel = slog.Level(-8)
)

var levels = map[slog.Leveler]string{
	TraceLevel: "TRACE",
}

func Init(level string, commit, version string) {
	logLevel := getLevel(level)
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: logLevel < slog.LevelInfo,
		AddSource: false,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levels[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}
			return a
		},
	})

	logger := slog.New(handler)

	if buildInfo, ok := debug.ReadBuildInfo(); ok && logLevel < slog.LevelDebug {
		logger = logger.With(
			slog.Group("buildInfo",
				slog.String("commit", commit),
				slog.String("verison", version),
				slog.String("go_version", buildInfo.GoVersion),
			),
		)
	}

	slog.SetDefault(logger)
	slog.Debug("configured logs")
}

func getLevel(level string) slog.Level {
	switch level {
	case debugLevel:
		return slog.LevelDebug
	case infoLevel:
		return slog.LevelInfo
	case warnLevel:
		return slog.LevelWarn
	case errorLevel:
		return slog.LevelError
	case traceLevel:
		return TraceLevel
	default:
		return slog.LevelDebug
	}
}

func Duration(start time.Time) slog.Attr {
	return slog.String(duration, time.Since(start).String())
}

func Trace(msg string, args ...any) {
	slog.Log(context.Background(), TraceLevel, msg, args...)
}

func Err(err error) slog.Attr {
	if err != nil {
		return slog.String(errStr, err.Error())
	}
	return slog.String(empty, "")
}
