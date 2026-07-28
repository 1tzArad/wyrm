package network

import (
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]*WSClient
	register   chan *WSClient
	unregister chan *WSClient
	broadcast  chan []byte

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*WSClient),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		broadcast:  make(chan []byte),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.registerClientHandler(client)
		case client := <-hub.unregister:
			hub.unregisterClientHandler(client)
		case message := <-hub.broadcast:
			hub.handleBroadcasts(message)
		}
	}
}

func (hub *Hub) registerClientHandler(client *WSClient) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.clients[client.UUID] = client
}

func (hub *Hub) unregisterClientHandler(client *WSClient) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if _, ok := hub.clients[client.UUID]; ok {
		delete(hub.clients, client.UUID)
		close(client.Send)
	}
}

func (hub *Hub) handleBroadcasts(message []byte) {
	for clientUUID, client := range hub.clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(hub.clients, clientUUID)
		}
	}
}

func (hub *Hub) Broadcast(message []byte) {
	hub.broadcast <- message
}

func (hub *Hub) Register(client *WSClient) {
	hub.register <- client
}
