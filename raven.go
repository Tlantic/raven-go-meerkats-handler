package raven_meerkats

import (
	"fmt"
	Meerkats "github.com/Tlantic/meerkats"
	Raven "github.com/getsentry/raven-go"
)

const (
	FIELD_TAGS = "Tags"
)

var levels = [...]Raven.Severity {
	Meerkats.LEVEL_TRACE: Raven.DEBUG,
	Meerkats.LEVEL_DEBUG: Raven.DEBUG,
	Meerkats.LEVEL_INFO: Raven.INFO,
	Meerkats.LEVEL_WARNING: Raven.WARNING,
	Meerkats.LEVEL_ERROR: Raven.ERROR,
	Meerkats.LEVEL_FATAL: Raven.FATAL,
	Meerkats.LEVEL_PANIC: Raven.FATAL,
}


type RavenHandler struct {
	*Raven.Client
}


//noinspection GoUnusedExportedFunction
func New(dsn string, tags map[string]string) (*RavenHandler, error) {

	if ( tags == nil ) {
		tags = make(map[string]string)
	}

	tags["logger"] = "meerkats"

	if client, err := Raven.NewClient(dsn, tags); err != nil {
		return nil, err
	} else {
		return &RavenHandler{client}, nil
	}

}


func (h *RavenHandler) HandleEntry(e Meerkats.Entry, done Meerkats.Callback) {
	defer done()

	p := Raven.NewPacket(e.Message)
	p.EventID = e.TraceId
	p.Timestamp = Raven.Timestamp(e.Timestamp)
	p.Level = levels[e.Level]



	tags := make(map[string]string)
	t := Meerkats.Fields{}
	if  e.Fields[FIELD_TAGS] != nil {
		t.Merge(e.Fields[FIELD_TAGS])
		for key, value := range t {
			tags[key] = fmt.Sprint(value)
		}
		delete(e.Fields, FIELD_TAGS)
	}


	for key, value := range e.Fields {
		p.Extra[key] = value
	}



	id, ch := h.Capture(p, tags);
	err := <- ch
	if err != nil {
		fmt.Errorf("%s: %s", id, err.Error())
	}
}