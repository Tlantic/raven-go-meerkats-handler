package raven_meerkats

import (
	"testing"
	"github.com/Tlantic/meerkats"
	"github.com/getsentry/raven-go"
)

var err = raven.SetDSN("https://5f0a7eb4fbf74558b5f4e01ca650bc72:f8413f5aac2e4370a550a31cf73d6489@sentry.io/115378")

func TestNew(t *testing.T) {
	if (err != nil) {
		t.Fatalf("%s", err.Error())
	}
	h := New()
	defer h.Dispose()
	logger := meerkats.New(meerkats.TRACE, h)
	logger.SetMeta("isMeta", "true")
	logger.With(meerkats.String("id", "12345"))
	logger.Error("hello", meerkats.Bool("test", true))
}
