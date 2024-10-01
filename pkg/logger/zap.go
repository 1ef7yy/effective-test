package logger

import (
	"go.uber.org/zap"
)

type log struct {
	*zap.Logger
}

func newZap(outputs []string) log {
	cfg := zap.NewDevelopmentConfig()

	cfg.Encoding = "json"
	if outputs != nil {
		cfg.OutputPaths = outputs
	}

	return log{
		zap.Must(cfg.Build()),
	}
}

func (l log) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l log) Info(msg string) {
	l.Logger.Info(msg)
}

func (l log) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l log) Error(msg string) {
	l.Logger.Error(msg)
}

func (l log) Fatal(msg string) {
	l.Logger.Fatal(msg)
}
