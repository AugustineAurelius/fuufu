package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerOpt func(c *zap.Config)

func WithDebug() LoggerOpt {
	return func(config *zap.Config) {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
}

func WithJSON() LoggerOpt {
	return func(config *zap.Config) {
		config.Encoding = "json"
	}
}

func New(opts ...LoggerOpt) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	for _, opt := range opts {
		opt(&config)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
