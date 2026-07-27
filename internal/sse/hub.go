package sse

import (
	"encoding/json"
	"log/slog"
	"sync"
)

// event sent to all sse clients
type Event struct {
	Type    string      `json:"type"`
	Webhook interface{} `json:"webhook"`
}

// hub manages sse client connections and event broadcasting
type Hub struct {
	mu      sync.RWMutex
	clients map[chan []byte]struct{}
}

// creates a new hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[chan []byte]struct{}),
	}
}

// adds a client and returns its channel
func (h *Hub) Register() chan []byte {
	ch := make(chan []byte, 64)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	slog.Info("SSE client connected", "total_clients", h.ClientCount())
	return ch
}

// removes a client and closes its channel
func (h *Hub) Unregister(ch chan []byte) {
	h.mu.Lock()
	if _, ok := h.clients[ch]; ok {
		delete(h.clients, ch)
		close(ch)
	}
	h.mu.Unlock()
	slog.Info("SSE client disconnected", "total_clients", h.ClientCount())
}

// sends an event to all clients, drops if buffer full
func (h *Hub) Broadcast(event *Event) {
	data, err := json.Marshal(event)
	if err != nil {
		slog.Error("failed to marshal SSE event", "error", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for ch := range h.clients {
		select {
		case ch <- data:
		default:
			slog.Warn("SSE client buffer full, dropping event")
		}
	}
}

// returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
