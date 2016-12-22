package raven_meerkats

import (
	"fmt"
	log "github.com/Tlantic/meerkats"
	"github.com/getsentry/raven-go"
	"sync"
	"time"
)


var levels = [...]raven.Severity {
	log.TRACE: raven.DEBUG,
	log.DEBUG: raven.DEBUG,
	log.INFO: raven.INFO,
	log.WARNING: raven.WARNING,
	log.ERROR: raven.ERROR,
	log.FATAL: raven.FATAL,
	log.PANIC: raven.FATAL,
}

var _ log.Handler = (*RavenHandler)(nil)
type RavenHandler struct {
	*raven.Client
	Level log.Level
	TimeLayout string
	mu sync.Mutex
	tags map[string]string
	fields map[string]interface{}
}




//noinspection GoUnusedExportedFunction
func New(options...log.HandlerOption) (h *RavenHandler) {
	h = &RavenHandler{
		Client: raven.DefaultClient,
		Level: log.LEVEL_ALL,
		TimeLayout: "",
		mu: sync.Mutex{},
		tags: map[string]string{},
		fields: map[string]interface{}{},
	}
	for _, opt := range options {
		opt.Apply(h)
	}
	return h
}

func Register(options ...log.HandlerOption) log.LoggerOption {
	return log.LoggerReceiver(func(l log.Logger) {
		l.Register(New(options...))
	})
}

func (h *RavenHandler) Apply(l log.Logger) {
	l.Register(h)
}

func (h *RavenHandler) AddBool(key string, value bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddString(key string, value string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddInt(key string, value int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddInt64(key string, value int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddUint(key string, value uint) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddUint64(key string, value uint64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddFloat32(key string, value float32) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) AddFloat64(key string, value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) Add(key string, value interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.fields[key] = value
}
func (h *RavenHandler) With(fields ...log.Field) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, f := range fields {
		h.fields[f.Key] = f.Get()
	}
}

func (h *RavenHandler) Log(t time.Time, level log.Level, msg string, fields []log.Field, meta map[string]string) {
	if ( h.Level&level != 0 ) {
		packet := raven.NewPacket(msg)
		packet.Level = levels[level]
		packet.AddTags(h.Tags)
		packet.AddTags(meta)
		packet.Logger = "meerkats"
		packet.Timestamp = raven.Timestamp(t)
		for k, v := range h.fields {
			packet.Extra[k] = v
		}
		for _, f := range fields {
			packet.Extra[f.Key] = f.Get()
		}

		id, ch := h.Capture(packet, h.tags);
		err := <- ch
		if err != nil {
			fmt.Errorf("%s: %s", id, err.Error())
		}
	}
}

func (h *RavenHandler) SetTimeLayout(layout string) {
	h.TimeLayout = layout
}
func (h *RavenHandler) GetTimeLayout() string {
	return h.TimeLayout
}

func (h *RavenHandler) SetLevel(level log.Level) {
	h.Level = level
}
func (h *RavenHandler) GetLevel() log.Level {
	return h.Level
}


func (h *RavenHandler) Clone() log.Handler {
	clone := &RavenHandler{
		Client: h.Client,
		Level: h.Level,
		TimeLayout: h.TimeLayout,
		tags: map[string]string{},
		fields: map[string]interface{}{},
		mu: sync.Mutex{},
	}
	for k, v := range h.tags {
		clone.tags[k] = v
	}
	for k, v := range h.fields {
		clone.fields[k] = v
	}
	return clone
}
func (h *RavenHandler) Dispose() {

}