package network

import "encoding/json"

type Message struct {
	Type      string          `json:"type"`
	RequestID string          `json:"request_id,omitempty"`
	Payload   json.RawMessage `json:"payload"`
}

type ServerErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
