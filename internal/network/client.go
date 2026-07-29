package network

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type WSClient struct {
	UUID       uuid.UUID
	Hub        *Hub
	Conn       *websocket.Conn
	send       chan []byte
	PlayerUUID uuid.UUID
	Registery  *Registery
}

func NewClient(conn *websocket.Conn, hub *Hub, playerUUID uuid.UUID, registery *Registery) *WSClient {
	return &WSClient{
		UUID:       uuid.New(),
		Hub:        hub,
		Conn:       conn,
		send:       make(chan []byte),
		PlayerUUID: playerUUID,
		Registery:  registery,
	}
}

// reads client's sends
func (client *WSClient) ReadPump() {
	defer func() {
		client.Hub.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(appData string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Infof("Client %s disconnected unexpectedly: %v", client.UUID, err)
			}
			break
		}
		log.Debugf("Received: %s", string(message))
		client.Registery.Dispatch(client, message)
	}
}

// write whatever in client's Send channel.
func (client *WSClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := client.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (client *WSClient) Send(message []byte) {
	client.send <- message
}
