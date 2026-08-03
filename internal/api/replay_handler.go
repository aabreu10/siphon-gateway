package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/aabreu10/siphon-gateway/internal/broker"
	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/aabreu10/siphon-gateway/internal/sse"
)

// handles POST /api/v1/webhook/{id}/replay — re-queues a failed webhook
func replayHandler(repo *database.WebhookRepo, pub *broker.Publisher, hub *sse.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			http.Error(w, `{"error":"invalid webhook ID"}`, http.StatusBadRequest)
			return
		}

		// fetch webhook
		webhook, err := repo.GetByID(r.Context(), id)
		if err != nil {
			slog.Error("failed to fetch webhook for replay", "error", err, "webhook_id", id)
			http.Error(w, `{"error":"webhook not found"}`, http.StatusNotFound)
			return
		}

		if webhook.Status != "FAILED_DLQ" {
			http.Error(w, `{"error":"webhook is not eligible for replay"}`, http.StatusConflict)
			return
		}

		// reset status
		if err := repo.UpdateStatus(r.Context(), id, "PENDING", 0); err != nil {
			slog.Error("failed to reset webhook status", "error", err, "webhook_id", id)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// re-publish to queue
		msg := &broker.Message{
			WebhookID:  webhook.ID,
			Payload:    webhook.Payload,
			TargetURL:  webhook.TargetURL,
			RetryCount: 0,
			Source:     webhook.Source,
		}
		if err := pub.PublishToProcess(r.Context(), msg); err != nil {
			slog.Error("failed to republish webhook", "error", err, "webhook_id", id)
			http.Error(w, `{"error":"failed to queue replay"}`, http.StatusInternalServerError)
			return
		}

		// broadcast sse event
		hub.Broadcast(&sse.Event{
			Type: "replayed",
			Webhook: map[string]interface{}{
				"id":          webhook.ID,
				"source":      webhook.Source,
				"payload":     webhook.Payload,
				"target_url":  webhook.TargetURL,
				"status":      "PENDING",
				"retry_count": 0,
			},
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":     id,
			"status": "replayed",
		})

		slog.Info("webhook replayed", "webhook_id", id)
	}
}
