package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// handles POST /api/v1/echo — default target endpoint for webhook delivery that returns 200 ok
func echoHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, maxPayloadSize))
		if err != nil {
			slog.Error("echo receiver: failed to read request body", "error", err)
			http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
			return
		}
		slog.Info("echo receiver: webhook delivered successfully",
			"content_length", len(body),
			"source", r.Header.Get("X-Webhook-Source"),
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"received": true,
			"status":   "ok",
		})
	}
}
