package stream

import "sync"

type Client struct {
	UserID string
	Ch     chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string][]*Client // user_id → connections
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string][]*Client),
	}
}

func (h *Hub) Add(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[userID] = append(h.clients[userID], c)
}

func (h *Hub) Remove(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	list := h.clients[userID]
	for i, client := range list {
		if client == c {
			h.clients[userID] = append(list[:i], list[i+1:]...)
			break
		}
	}
}

func (h *Hub) Publish(userID string, data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, c := range h.clients[userID] {
		select {
		case c.Ch <- data:
		default:
		}
	}
}
