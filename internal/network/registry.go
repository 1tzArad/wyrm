package network

import (
	"encoding/json"

	"github.com/charmbracelet/log"
)

type HandlerFunc func(c *Context)

type Registery struct {
	handlers map[string]HandlerFunc
}

func NewRegistery() *Registery {
	return &Registery{handlers: make(map[string]HandlerFunc)}
}

func (r *Registery) Register(msgType string, handler HandlerFunc) {
	r.handlers[msgType] = handler
}

func (r *Registery) Dispatch(client *WSClient, raw []byte) {
	log.Debugf("Raw message received: %s", string(raw)) // این خط رو اضافه کن
	var msg Message
	if err := json.Unmarshal(raw, &msg); err != nil {
		log.Errorf("Invalid message from %s: %v", client.UUID, err)
		return
	}

	handler, ok := r.handlers[msg.Type]
	if !ok {
		log.Errorf("no handler registered for type: %s", msg.Type)
		return
	}

	ctx := &Context{
		ClientUUID: client.UUID,
		PlayerUUID: client.PlayerUUID,
		Payload:    msg.Payload,
		client:     client,
		hub:        client.Hub,
	}

	handler(ctx)
}
