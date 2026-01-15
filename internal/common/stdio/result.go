package stdio

import (
	"encoding/json"
	"fmt"
)

// HandlerResult encapsule le résultat d'un handler (succès ou échec).
type HandlerResult struct {
	successType string
	failType    string
	payload     any
	err         error
}

// Publish envoie le message approprié (succès ou échec) via le publisher.
func (r HandlerResult) Publish(msg Message, pub Publisher, eventTopic string) error {
	if r.err == nil {
		return pub.Publish(*msg.Reply(r.successType, r.payload, eventTopic))
	}

	data := map[string]any{"error": r.err.Error()}
	if extra, ok := r.payload.(map[string]any); ok {
		for k, v := range extra {
			data[k] = v
		}
	}
	return pub.Publish(*msg.Reply(r.failType, data, eventTopic))

}

// Success crée un résultat de succès.
func Success(eventType string, payload any) HandlerResult {
	return HandlerResult{successType: eventType, payload: payload}
}

// Fail crée un résultat d'échec.
func Fail(eventType string, err error, extra map[string]any) HandlerResult {
	return HandlerResult{failType: eventType, err: err, payload: extra}
}

// UnmarshalPayload décode le payload JSON dans dest.
func UnmarshalPayload(msg Message, dest any) error {
	if err := json.Unmarshal(msg.Payload, dest); err != nil {
		return fmt.Errorf("bad payload: %w", err)
	}
	return nil
}
