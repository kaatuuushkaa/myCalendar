package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(isDev bool) (*zap.Logger, error) {
	if isDev {
		return zap.NewDevelopment()
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}

func MustNew(isDev bool) *zap.Logger {
	log, err := New(isDev)
	if err != nil {
		panic(err)
	}
	return log
}
