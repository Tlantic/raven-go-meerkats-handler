package raven_meerkats

import (
	"testing"
	"github.com/Tlantic/meerkats"
	"github.com/getsentry/raven-go"
)

var err = raven.SetDSN("")

var TAGS = map[string]string{
	"test": "true",
}

func TestNew(t *testing.T) {
	if (err != nil) {
		t.Fatalf("%s", err.Error())
	}

	logger := meerkats.New(meerkats.TRACE, Register())
	logger.SetMeta("isMeta", "true")
	logger.With(meerkats.String("id", "12345"))
	logger.Error("hello", meerkats.Bool("test", true))
}
