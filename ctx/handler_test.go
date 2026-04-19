package ctx

import (
	"bytes"
	"log/slog"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/require"
)

func getLogger() (*slog.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	json := slog.NewJSONHandler(buf, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
		// ReplaceAttr: logger.ReplaceAttr,
	})

	h := NewAttrsHandler(json, "ctx_")

	return slog.New(h), buf
}

func TestCtxAttr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log, buf := getLogger()
		ctx := WithLogAttrs(t.Context())

		attrs := ExtractAttrs(ctx)
		attrs.Set(
			slog.Bool("is_test", true),
			slog.String("key", "value"),
		)

		log.WarnContext(ctx, "warning", "is_test", false)

		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","level":"WARN","msg":"warning","is_test":false,"ctx_is_test":true,"ctx_key":"value"}`,
			buf.String(),
		)
	})
}

func TestWithAttrs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log, buf := getLogger()

		ctx := WithLogAttrs(t.Context(), slog.Bool("is_test", true))

		log.WarnContext(ctx, "warning")
		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","level":"WARN","msg":"warning","ctx_is_test":true}`,
			buf.String(),
		)
		buf.Reset()

		log.WarnContext(ctx, "warning", slog.Bool("is_test", false))
		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","is_test":false,"level":"WARN","msg":"warning","ctx_is_test":true}`,
			buf.String(),
		)
		buf.Reset()

		ctx2 := WithLogAttrs(ctx, slog.Bool("child", true))
		log.WarnContext(ctx2, "warning")
		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","level":"WARN","msg":"warning","ctx_is_test":true,"ctx_child":true}`,
			buf.String(),
		)
		buf.Reset()

		log.WarnContext(ctx, "warning")
		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","level":"WARN","msg":"warning","ctx_is_test":true}`,
			buf.String(),
		)
	})
}

func TestUnsetAttrs(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log, buf := getLogger()
		ctx := WithLogAttrs(t.Context())

		attrs := ExtractAttrs(ctx)
		attrs.Set(
			slog.Bool("is_test", true),
			slog.String("key", "value"),
		)
		attrs.Unset("key")

		log.WarnContext(ctx, "warning")

		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","level":"WARN","msg":"warning","ctx_is_test":true}`,
			buf.String(),
		)
	})
}

func TestChildHandler(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		log, buf := getLogger()
		log = log.With(slog.Bool("sub-logger", true))
		log = log.WithGroup("group")

		log.WarnContext(t.Context(), "warning", slog.Int("cpu", 12))
		require.JSONEq(t,
			`{"time":"2000-01-01T03:00:00+03:00","level":"WARN","msg":"warning","sub-logger":true,"group":{"cpu":12}}`,
			buf.String(),
		)
	})
}
