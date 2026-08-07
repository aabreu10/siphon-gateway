package api

import (
	"encoding/json"
	"io"
	"log/slog"
	"mime"
	"net/http"

	"github.com/aabreu10/siphon-gateway/internal/broker"
	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/aabreu10/siphon-gateway/internal/sse"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const maxPayloadSize = 1 << 20 // 1 MB

// handles POST /api/v1/ingest/{endpoint_id} — validates, persists, publishes, returns id
func ingestHandler(repo *database.WebhookRepo, pub *broker.Publisher, hub *sse.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "*")
		
		var endpointID uuid.UUID
		var targetURL string
		var secretKey string

		if param != "" {
			// Try to parse the wildcard parameter as a UUID
			if id, err := uuid.Parse(param); err == nil {
				endpointID = id
			} else {
				// It's not a UUID, so treat it as a direct target URL (Simulator mode)
				targetURL = param
			}
		}

		// check content-type
		ct := r.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(ct)
		if err != nil || mediaType != "application/json" {
			http.Error(w, `{"error":"Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
			return
		}

		// read body with size limit
		body, err := io.ReadAll(io.LimitReader(r.Body, maxPayloadSize))
		if err != nil {
			slog.Error("failed to read request body", "error", err)
			http.Error(w, `{"error":"failed to read body"}`, http.StatusBadRequest)
			return
		}

		// parse json
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, `{"error":"invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		// get source from header or query
		source := r.Header.Get("X-Webhook-Source")
		if source == "" {
			source = r.URL.Query().Get("source")
		}
		if source == "" {
			source = "unknown"
		}

		if endpointID != uuid.Nil {
			// Real endpoint mode
			endpoint, err := repo.GetEndpointForIngest(r.Context(), endpointID)
			if err != nil {
				slog.Error("endpoint not found", "error", err, "endpoint_id", endpointID)
				http.Error(w, `{"error":"endpoint not found"}`, http.StatusNotFound)
				return
			}
			targetURL = endpoint.TargetURL
			secretKey = endpoint.SecretKey
		} else {
			// Simulator mode
			if targetURL == "" {
				targetURL = r.URL.Query().Get("target_url")
			}
		}

		// save to db
		id, err := repo.Insert(r.Context(), source, payload, targetURL, endpointID)
		if err != nil {
			slog.Error("failed to insert webhook", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// publish to queue
		msg := &broker.Message{
			WebhookID:  id,
			EndpointID: endpointID,
			Payload:    payload,
			TargetURL:  targetURL,
			SecretKey:  secretKey,
			RetryCount: 0,
			Source:     source,
		}
		if err := pub.PublishToProcess(r.Context(), msg); err != nil {
			slog.Error("failed to publish to queue", "error", err, "webhook_id", id)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// broadcast sse event
		hub.Broadcast(&sse.Event{
			Type: "received",
			Webhook: map[string]interface{}{
				"id":          id,
				"source":      source,
				"payload":     payload,
				"target_url":  targetURL,
				"status":      "PENDING",
				"retry_count": 0,
			},
		})

		// respond with id
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     id,
			"status": "received",
		})

		slog.Info("webhook ingested",
			"webhook_id", id,
			"source", source,
		)
	}
}
