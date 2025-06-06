package log

import (
	"log/slog"
	"os"
	"time"
)

const (
	duration = "duration"
	empty    = "empty"
	errStr   = "error"

	debugLevel string = "debug"
	infoLevel  string = "info"
	warnLevel  string = "warn"
	errorLevel string = "error"
	traceLevel string = "trace"

	Trace = slog.Level(-8)
)

var levels = map[slog.Leveler]string{
	Trace: "TRACE",
}

func Init(level string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     getLevel(level),
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
	slog.SetDefault(slog.New(handler))
	slog.Info("configured logs")
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
		return Trace
	default:
		return slog.LevelDebug
	}
}

func Duration(start time.Time) slog.Attr {
	return slog.String(duration, time.Since(start).String())
}

func Err(err error) slog.Attr {
	if err != nil {
		return slog.String(errStr, err.Error())
	}
	return slog.String(empty, "")
}
