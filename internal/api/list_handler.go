package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// handles GET /api/v1/webhooks — paginated list
func listHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		offset := 0

		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 200 {
				limit = v
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil && v >= 0 {
				offset = v
			}
		}
		
		status := r.URL.Query().Get("status")
		if status == "ALL" {
			status = ""
		}
		search := r.URL.Query().Get("search")

		userIDStr, ok := r.Context().Value(UserIDKey).(string)
		if !ok {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		webhooks, total, err := repo.ListRecent(r.Context(), limit, offset, status, search, userID)
		if err != nil {
			slog.Error("failed to list webhooks", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// return empty array instead of null
		if webhooks == nil {
			webhooks = []database.Webhook{}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Total-Count", strconv.Itoa(total))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"webhooks": webhooks,
			"total":    total,
			"limit":    limit,
			"offset":   offset,
		})
	}
}

// handles GET /api/v1/webhook/{id}/logs
func logsHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid id format"}`, http.StatusBadRequest)
			return
		}

		logs, err := repo.GetDeliveryLogs(r.Context(), id)
		if err != nil {
			slog.Error("failed to get delivery logs", "error", err, "id", id)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		if logs == nil {
			logs = []database.DeliveryLog{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}
