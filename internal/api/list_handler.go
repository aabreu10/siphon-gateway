package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aabreu10/siphon-gateway/internal/database"
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

		webhooks, total, err := repo.ListRecent(r.Context(), limit, offset)
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
