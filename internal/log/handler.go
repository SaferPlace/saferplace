package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
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
	builder := new(strings.Builder)
	if !r.Time.IsZero() {
		builder.WriteString(r.Time.Format(time.RFC3339Nano) + " ")
	}
	builder.WriteString(r.Level.String() + "\t" + r.Message + "\t")
	r.Attrs(func(a slog.Attr) bool {
		h.appendAttr(builder, a)
		return true
	})

	extraAttr := h.attrs[:]
	// Add the tracing information if available.
	span := trace.SpanFromContext(ctx)
	if traceID := span.SpanContext().TraceID(); traceID.IsValid() {
		extraAttr = append(extraAttr,
			slog.String("trace_id", traceID.String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	for _, attr := range extraAttr {
		h.appendAttr(builder, attr)
	}

	fmt.Fprintln(h.out, builder.String())
	return nil
}

func (h *handler) appendAttr(builder *strings.Builder, a slog.Attr) {
	a.Value = a.Value.Resolve()
	// ignore empty attr
	if a.Equal(slog.Attr{}) {
		return
	}
	switch a.Value.Kind() {
	case slog.KindString:
		builder.WriteString(fmt.Sprintf("%s=%q", a.Key, a.Value.String()) + " ")
	case slog.KindTime:
		builder.WriteString(a.Key + "=" + a.Value.Time().Format(time.RFC3339Nano) + " ")
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return
		}
		for _, ga := range attrs {
			ga.Key = a.Key + "." + ga.Key
			h.appendAttr(builder, ga)
		}
	default:
		builder.WriteString(fmt.Sprintf("%s=%s", a.Key, a.Value.String()) + " ")
	}
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
