package internal

import (
	"context"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gface"
	"go.uber.org/zap"
)

type Logger struct {
	gface.ILogger
	module string
	Zap    *zap.SugaredLogger
}

func NewLogger(module string, zapLogger *zap.SugaredLogger) *Logger {
	return &Logger{
		ILogger: gface.NewLogger(module, zapLogger),
		module:  module,
		Zap:     zapLogger,
	}
}

func (l *Logger) DebugKVs(ctx context.Context, msg string, kv ...any) {
	if ctx != nil {
		l.Zap.Debug(gcommon.GetRequestIdFormat(ctx) + gcommon.Kv2Str(msg, kv...))
	} else {
		l.Zap.Debug(gcommon.Kv2Str(msg, kv...))
	}
}

func (l *Logger) InfoKVs(ctx context.Context, msg string, kv ...any) {

	if ctx != nil {
		l.Zap.Info(gcommon.GetRequestIdFormat(ctx) + gcommon.Kv2Str(msg, kv...))
	} else {
		l.Zap.Info(gcommon.Kv2Str(msg, kv...))
	}
}

func (l *Logger) WarnKVs(ctx context.Context, msg string, kv ...any) {
	if ctx != nil {
		l.Zap.Warn(gcommon.GetRequestIdFormat(ctx) + gcommon.Kv2Str(msg, kv...))
	} else {
		l.Zap.Warn(gcommon.Kv2Str(msg, kv...))
	}
}

func (l *Logger) ErrorKVs(ctx context.Context, msg string, kv ...any) {
	if ctx != nil {
		l.Zap.Error(gcommon.GetRequestIdFormat(ctx) + gcommon.Kv2Str(msg, kv...))
	} else {
		l.Zap.Error(gcommon.Kv2Str(msg, kv...))
	}
}

func (l *Logger) PanicKVs(ctx context.Context, msg string, kv ...any) {
	if ctx != nil {
		l.Zap.Panic(gcommon.GetRequestIdFormat(ctx) + gcommon.Kv2Str(msg, kv...))
	} else {
		l.Zap.Panic(gcommon.Kv2Str(msg, kv...))
	}
}
