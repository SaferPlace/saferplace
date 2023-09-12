// Package log is a helper library for [log/slog] which exposes a nicer, fully context based API.
package log

import (
	"context"
	"log/slog"
)

// Logger is an abstracted interface over the [log/slog] packages which provides only a context
// based API.
type Logger interface {
	Debug(context.Context, string, ...slog.Attr)
	Info(context.Context, string, ...slog.Attr)
	Warn(context.Context, string, ...slog.Attr)
	Error(context.Context, string, ...slog.Attr)

	With(...slog.Attr) Logger
}

var _ Logger = (*logger)(nil)

type logger struct {
	l *slog.Logger
}

type Handler = slog.Handler

// New creates a new logger
func New(h Handler) *logger {
	return &logger{
		l: slog.New(h),
	}
}

// Debug Level Logging
func (l *logger) Debug(ctx context.Context, msg string, args ...slog.Attr) {
	l.l.LogAttrs(ctx, slog.LevelDebug, msg, args...)
}

// Info Level Logging
func (l *logger) Info(ctx context.Context, msg string, args ...slog.Attr) {
	l.l.LogAttrs(ctx, slog.LevelInfo, msg, args...)
}

// Warn Level Logging
func (l *logger) Warn(ctx context.Context, msg string, args ...slog.Attr) {
	l.l.LogAttrs(ctx, slog.LevelDebug, msg, args...)
}

// Error Level Logging
func (l *logger) Error(ctx context.Context, msg string, args ...slog.Attr) {
	l.l.LogAttrs(ctx, slog.LevelInfo, msg, args...)
}

// With creates a new logger with added args.
func (l *logger) With(args ...slog.Attr) Logger {
	if len(args) == 0 {
		return l
	}
	return New(l.l.Handler().WithAttrs(args))
}

// Error is the standard Entry for when we want to log an error
func Error(err error) slog.Attr {
	return slog.String("error", err.Error())
}
