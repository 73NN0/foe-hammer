package stdio

import (
	"encoding/json"
	"io"
	"sync"
)

// Publisher définit l'interface pour publier des messages.
type Publisher interface {
	Publish(msg Message) error
}

// StdoutPublisher publie des messages JSON sur un io.Writer.
type StdoutPublisher struct {
	w   io.Writer
	enc *json.Encoder
	mu  sync.Mutex
}

// NewStdoutPublisher crée un publisher qui écrit sur w.
func NewStdoutPublisher(w io.Writer) *StdoutPublisher {
	return &StdoutPublisher{
		w:   w,
		enc: json.NewEncoder(w),
	}
}

// Publish encode et écrit le message (thread-safe).
func (p *StdoutPublisher) Publish(msg Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.enc.Encode(msg)
}
