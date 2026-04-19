package ctx

import (
	"context"
	"log/slog"
)

// AttrsHandler — декоратор Handler для обработки атрибутов логов из контекста
type AttrsHandler struct {
	slog.Handler

	prefix string
}

func NewAttrsHandler(slog slog.Handler, prefix string) slog.Handler {
	return &AttrsHandler{
		Handler: slog,
		prefix:  prefix,
	}
}

func (h *AttrsHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := ExtractAttrs(ctx)
	for _, attr := range attrs {
		key := h.prefix + attr.Key
		record.AddAttrs(slog.Attr{Key: key, Value: attr.Value})
	}

	// делегируем обработку
	return h.Handler.Handle(ctx, record)
}

func (h *AttrsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &AttrsHandler{
		Handler: h.Handler.WithAttrs(attrs),
		prefix:  h.prefix,
	}
}

func (h *AttrsHandler) WithGroup(name string) slog.Handler {
	return &AttrsHandler{
		Handler: h.Handler.WithGroup(name),
		prefix:  h.prefix,
	}
}
