package logger

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps a zap.Logger so the rest of the codebase depends on a small,
// stable interface rather than zap directly.
type Logger struct {
	*zap.Logger
}

// New constructs a Logger from level + format strings.
// level: "debug" | "info" | "warn" | "error"
// format: "console" | "json"
func New(level, format string) (*Logger, error) {
	lvl, err := zapcore.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zapcore.InfoLevel
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if strings.EqualFold(format, "json") {
		encoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		lvl,
	)

	// Disable stacktraces for warn and below to keep logs readable.
	l := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(0))
	return &Logger{Logger: l}, nil
}

// Sync flushes buffered log entries. Safe to call at shutdown.
func (l *Logger) Sync() {
	if l == nil || l.Logger == nil {
		return
	}
	_ = l.Logger.Sync()
}

// Must panics if the logger cannot be constructed.
func Must(level, format string) *Logger {
	l, err := New(level, format)
	if err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}
	return l
}
