package network

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type Hub struct {
	clients        map[uuid.UUID]*WSClient
	playerToClient map[uuid.UUID]*WSClient
	register       chan *WSClient
	unregister     chan *WSClient
	broadcast      chan []byte

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[uuid.UUID]*WSClient),
		playerToClient: make(map[uuid.UUID]*WSClient),
		register:       make(chan *WSClient),
		unregister:     make(chan *WSClient),
		broadcast:      make(chan []byte),
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
	hub.playerToClient[client.PlayerUUID] = client
	log.Debug("players to client")
	log.Debug(hub.playerToClient)
	log.Debug("clients")
	log.Debug(hub.clients)
}

func (hub *Hub) unregisterClientHandler(client *WSClient) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if _, ok := hub.clients[client.UUID]; ok {
		delete(hub.clients, client.UUID)
		delete(hub.playerToClient, client.PlayerUUID)
		close(client.Send)
	}
}

func (hub *Hub) handleBroadcasts(message []byte) {
	for _, client := range hub.clients {
		select {
		case client.Send <- message:
		default:
			hub.unregisterClientHandler(client)
		}
	}
}

func (hub *Hub) Broadcast(message []byte) {
	hub.broadcast <- message
}

func (hub *Hub) Register(client *WSClient) {
	hub.register <- client
}

func (hub *Hub) GetPlayer(uuid uuid.UUID) *WSClient {
	p, ok := hub.clients[uuid]
	if !ok {
		return nil
	}
	return p
}

func (hub *Hub) SendToPlayer(playerUUID uuid.UUID, message []byte) error {
	client, ok := hub.playerToClient[playerUUID]
	if !ok {
		return fmt.Errorf("player %s not connected", playerUUID)
	}

	select {
	case client.Send <- message:
		return nil
	default:
		return fmt.Errorf("failed to send to player %s: channel full", playerUUID)
	}
}
