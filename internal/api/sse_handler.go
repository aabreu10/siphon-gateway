package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aabreu10/siphon-gateway/internal/sse"
)

// handles GET /api/v1/events — streams real-time webhook updates
func sseHandler(hub *sse.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check flusher support
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		// sse headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no") // disable nginx buffering

		// register client
		clientCh := hub.Register()
		defer hub.Unregister(clientCh)

		// send connected event
		fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
		flusher.Flush()

		// stream until disconnect
		for {
			select {
			case <-r.Context().Done():
				slog.Info("SSE client disconnected")
				return
			case data, ok := <-clientCh:
				if !ok {
					return
				}
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}
