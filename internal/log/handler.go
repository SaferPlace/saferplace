package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"
)

// handler is a new handler with a more human readable output which takes less space.
type handler struct {
	out   io.Writer
	level slog.Level

	attrs  []slog.Attr
	groups []string
}

var _ slog.Handler = (*handler)(nil)

func (h *handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func NewHandler(out io.Writer, level slog.Level) *handler {
	return &handler{
		out:   out,
		level: level,
	}
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	kvs := make([]string, 0, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		kvs = append(kvs, strings.Join(append(h.groups, a.Key), ".")+"="+a.Value.String())
		return true
	})
	fmt.Fprintf(h.out, "%s %s\t%s\t%s\n",
		r.Time.Format(time.RFC3339Nano), r.Level, r.Message, strings.Join(kvs, " "))
	return nil
}

func (h *handler) clone() *handler {
	c := *h
	return &c
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	c := h.clone()
	c.attrs = append(c.attrs, attrs...)
	return c
}

func (h *handler) WithGroup(group string) slog.Handler {
	c := h.clone()
	c.groups = append(c.groups, group)
	return c
}
