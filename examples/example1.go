package main

import (
	"fmt"
	Meerkats "github.com/Tlantic/meerkats"
	Raven "github.com/getsentry/raven-go"
	"reflect"
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


	e.Fields["author"] = "João"
	tags := make(map[string]string)
	if e.Fields[FIELD_TAGS] != nil {

		rElem := reflect.ValueOf(e.Fields[FIELD_TAGS]).Elem()
		switch rElem.Kind() {
		case reflect.Struct:
			for i := 0; i < rElem.NumField(); i++ {
				valueField := rElem.Field(i)
				typeField := rElem.Type().Field(i)
				tags[typeField.Name] = fmt.Sprint(valueField)
			}
		case reflect.Map:
			for _, v := range rElem.MapKeys() {
				if k, ok := v.Interface().(string); ok {
					tags[k] = fmt.Sprint(rElem.MapIndex(v))
				}
			}
		}
		delete(e.Fields, FIELD_TAGS)
	}


	for key, value := range e.Fields {
		if v, ok := value.(reflect.Value); ok {
			p.Extra[key] = v.Interface()
		} else {
			p.Extra[key] = value
		}

	}



	id, ch := h.Capture(p, tags);
	err := <- ch
	if err != nil {
		fmt.Errorf("%s: %s", id, err.Error())
	}
}
const (
	DSN = ""
)

var TAGS = map[string]string{
	"test": "true",
}
func main() {


	r , err := New(DSN, TAGS)
	if err != nil {
		panic(r)
	}

	m := Meerkats.New(nil)
	defer func() {
		if r := recover(); r != nil {
			m.Panic(r)
		}
		m.Close()
	}()
	m.RegisterHandler(Meerkats.LEVEL_ALL, r.HandleEntry)

	entry := m.WithFields(Meerkats.Fields{
		"Field1": true,
		"Field2": "true",
		"Field3": nil,
		"Field4": false,
		FIELD_TAGS: Meerkats.Fields{
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