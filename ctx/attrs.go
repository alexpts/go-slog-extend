package ctx

import (
	"log/slog"
)

// Attrs
// @thread unsafe
type Attrs map[string]slog.Attr

func NewAttrs(capacity int) Attrs {
	return make(map[string]slog.Attr, capacity)
}

func (a Attrs) With(attrs ...slog.Attr) Attrs {
	newAttrs := make(map[string]slog.Attr, len(a)+len(attrs))

	for _, attr := range a {
		newAttrs[attr.Key] = attr
	}

	for _, attr := range attrs {
		newAttrs[attr.Key] = attr
	}

	return newAttrs
}

func (a Attrs) Set(attrs ...slog.Attr) {
	for _, attr := range attrs {
		a[attr.Key] = attr
	}
}

func (a Attrs) Unset(keys ...string) {
	for _, key := range keys {
		delete(a, key)
	}
}
