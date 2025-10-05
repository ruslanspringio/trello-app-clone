package ws

import (
	"log"
	"sync"
)

type Client struct {
	hub  *Hub
	conn Conn
	send chan []byte
}

type Hub struct {
	clients    map[int]map[*Client]bool
	broadcast  chan broadcastMessage
	register   chan *subscription
	unregister chan *subscription
	mu         sync.Mutex
}

type broadcastMessage struct {
	boardID int
	message []byte
}

type subscription struct {
	client  *Client
	boardID int
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int]map[*Client]bool),
		broadcast:  make(chan broadcastMessage),
		register:   make(chan *subscription),
		unregister: make(chan *subscription),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[sub.boardID]; !ok {
				h.clients[sub.boardID] = make(map[*Client]bool)
			}
			h.clients[sub.boardID][sub.client] = true
			h.mu.Unlock()
			log.Printf("Client registered to board %d", sub.boardID)

		case sub := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[sub.boardID]; ok {
				if _, ok := clients[sub.client]; ok {
					delete(clients, sub.client)
					close(sub.client.send)
					if len(clients) == 0 {
						delete(h.clients, sub.boardID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered from board %d", sub.boardID)

		case msg := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.clients[msg.boardID]; ok {
				for client := range clients {
					select {
					case client.send <- msg.message:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) BroadcastToBoard(boardID int, message []byte) {
	h.broadcast <- broadcastMessage{
		boardID: boardID,
		message: message,
	}
}
