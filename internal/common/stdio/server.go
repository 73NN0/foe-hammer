package stdio

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
)

type MessageHandler func(msg Message, pub Publisher) error

type ServerConfig struct {
	Topic   string
	Handler MessageHandler
}

type Server struct {
	config ServerConfig
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{config: cfg}
}

func (s *Server) Serve(r io.Reader, w io.Writer) error {
	dec := json.NewDecoder(bufio.NewReader(r))
	pub := NewStdoutPublisher(w)
	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		if msg.Topic != s.config.Topic {
			continue
		}
		if err := s.config.Handler(msg, pub); err != nil {
			_ = pub.Publish((*msg.ReplySameTopic(msg.Type+"Failed", map[string]any{
				"error":        err.Error(),
				"command_type": msg.Type,
			})))
		}
	}
}
