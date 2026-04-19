package ctx

import (
	"context"
	"log/slog"
)

type slogAttrKey struct{}

// ExtractAttrs - извлекает контейнер с атрибутами, его можно модифицировать
// Можно достать из ближайшего контекста в цепочке ctx и мутировать атрибуты, если иммутабельный вариант не подходит
func ExtractAttrs(ctx context.Context) Attrs {
	if attrs, isOk := ctx.Value(slogAttrKey{}).(Attrs); isOk {
		return attrs
	}

	return NewAttrs(0)
}

// WithLogAttrs - создает расширенную копию атрибутов и новый контекст с этими атрибутами
// @Immutable
func WithLogAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	return context.WithValue(ctx, slogAttrKey{}, ExtractAttrs(ctx).With(attrs...))
}
