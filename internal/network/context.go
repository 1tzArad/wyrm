package network

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Context struct {
	ClientUUID uuid.UUID
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
		Type:    msgType,
		Payload: data,
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	c.client.Send <- raw

	return nil
}
