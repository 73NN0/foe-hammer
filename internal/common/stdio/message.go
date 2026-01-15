package stdio

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
)

type Message struct {
	MessageID     string          `json:"message_id"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	CausationID   string          `json:"causation_id,omitempty"`
	Topic         string          `json:"topic"`
	Type          string          `json:"type"`
	Payload       json.RawMessage `json:"payload,omitempty"`
}

// Reply crée un message en réponse, avec un topic personnalisé.
func (m *Message) Reply(eventType string, payload any, topic string) *Message {
	return &Message{
		MessageID:     randID(),
		CorrelationID: m.CorrelationID,
		CausationID:   m.MessageID,
		Topic:         topic,
		Type:          eventType,
		Payload:       mustJSON(payload),
	}
}

// ReplySameTopic crée un message en réponse sur le même topic.
func (m *Message) ReplySameTopic(eventType string, payload any) *Message {
	return m.Reply(eventType, payload, m.Topic)
}

// NewMessage crée un nouveau message (pas une réponse).
func NewMessage(eventType, topic string, payload any) *Message {
	return &Message{
		MessageID:     randID(),
		CorrelationID: randID(),
		Topic:         topic,
		Type:          eventType,
		Payload:       mustJSON(payload),
	}
}

func randID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
