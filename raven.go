package raven_go_meerkats_handler

import (
	"strings"
	"github.com/getsentry/raven-go"

	. "github.com/Tlantic/meerkats"
	"fmt"
)

const (
	FIELD_TAGS = "tags"
)

var levels = [...]string {
	LEVEL_TRACE: strings.ToLower(LEVEL_TRACE.String()),
	LEVEL_DEBUG: raven.DEBUG,
	LEVEL_INFO: raven.INFO,
	LEVEL_WARNING: raven.WARNING,
	LEVEL_ERROR: raven.ERROR,
	LEVEL_FATAL: raven.FATAL,
	LEVEL_PANIC: strings.ToLower(LEVEL_PANIC.String()),
}

type RavenHandler struct {
	raven.Client
}

//noinspection GoUnusedExportedFunction
func New(dsn string, tags map[string]string) (*RavenHandler, error) {

	if ( tags == nil ) {
		tags = make(map[string]string)
	}

	tags["Meerkats"] = true

	if client, err := raven.NewClient(dsn, tags); err != nil {
		return nil, err
	} else {
		return &RavenHandler{client}, nil
	}

}

func (h *RavenHandler) HandleEntry(e Entry, done Callback) {
	defer done()

	var tags map[string]string

	p := raven.NewPacket(e.Message)
	p.EventID = e.TraceId
	p.Timestamp = e.Timestamp
	p.Level = levels[e.Level]

	if value, ok := e.Fields[FIELD_TAGS].(map[string]string); ok {
		tags = value
		delete( e.Fields, FIELD_TAGS )
	}

	id, ch := h.Capture(p, tags);
	err := <- ch
	if err != nil {
		fmt.Errorf("%s: %s", id, err.Error())
	}
}