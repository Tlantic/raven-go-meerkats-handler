package raven_meerkats

import (
	"testing"
	"github.com/Tlantic/meerkats"
)

const (
	DSN = "https://5f0a7eb4fbf74558b5f4e01ca650bc72:f8413f5aac2e4370a550a31cf73d6489@sentry.io/115378"
)

var TAGS = map[string]string{
	"test": "true",
}

func TestNew(t *testing.T) {
	r , err := New(DSN, TAGS)
	if err != nil {
		t.Fatal(r)
	}

	// assert interface
	var _ meerkats.EntryHandler = r.HandleEntry
}

func TestRavenHandler_HandleEntry(t *testing.T) {

	r , err := New(DSN, TAGS)
	if err != nil {
		t.Fatal(r)
	}

	m := meerkats.New(nil)
	defer func() {
		if r := recover(); r != nil {
			m.Panic(r)
		}
		m.Close()
	}()
	m.RegisterHandler(meerkats.LEVEL_ALL, r.HandleEntry)

	entry := m.WithFields(meerkats.Fields{
		"Field1": true,
		"Field2": "true",
		"Field3": nil,
		"Field4": false,
		FIELD_TAGS: meerkats.Fields{
			"Tag1": "false",
			"Tag2": "true",
			"Tag3": "nil",
			"Tag4": "null",
		},
	})
/*
	entry.Trace("Trace message")
	entry.Debug("Debug message")
	entry.Info("Info message")
	entry.Warning("Warning message")
	entry.Error("Error message")*/
	entry.Fatal("Fatal message")
}