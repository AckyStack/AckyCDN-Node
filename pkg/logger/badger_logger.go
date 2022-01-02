package logger

import "github.com/gookit/slog"

type DBLogger struct{}

func (*DBLogger) Errorf(s string, i ...interface{}) {
	slog.Errorf(s, i)
}

func (*DBLogger) Warningf(s string, i ...interface{}) {
	slog.Warnf(s, i)
}

func (*DBLogger) Infof(s string, i ...interface{}) {
	slog.Infof(s, i)
}

func (*DBLogger) Debugf(s string, i ...interface{}) {
	slog.Debugf(s, i)
}
