package network

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Context struct {
	RequestID  string
	ClientUUID uuid.UUID
	PlayerUUID uuid.UUID
	Payload    json.RawMessage
	hub        *Hub
	client     *WSClient
}

func (c *Context) Bind(v interface{}) error {
	return json.Unmarshal(c.Payload, v)
}

func (c *Context) Reply(msgType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := Message{
		Type:      msgType,
		RequestID: c.RequestID,
		Payload:   data,
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.client.Send(raw)

	return nil
}

func (c *Context) ReplyError(code, message string) error {
	payload := ServerErrorPayload{
		Code:    code,
		Message: message,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := Message{
		Type:      "error",
		RequestID: c.RequestID,
		Payload:   data,
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.client.Send(raw)

	return nil
}
