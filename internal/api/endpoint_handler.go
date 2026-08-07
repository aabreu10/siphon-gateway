package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateEndpointRequest struct {
	Name      string `json:"name"`
	TargetURL string `json:"target_url"`
	SecretKey string `json:"secret_key"`
}

func createEndpointHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateEndpointRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.TargetURL == "" || req.SecretKey == "" {
			http.Error(w, `{"error":"name, target_url, and secret_key are required"}`, http.StatusBadRequest)
			return
		}

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

		id, err := repo.CreateEndpoint(r.Context(), req.Name, req.TargetURL, req.SecretKey, userID)
		if err != nil {
			slog.Error("failed to create endpoint", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         id,
			"name":       req.Name,
			"target_url": req.TargetURL,
		})
	}
}

func listEndpointsHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		endpoints, err := repo.ListEndpoints(r.Context(), userID)
		if err != nil {
			slog.Error("failed to list endpoints", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		if endpoints == nil {
			endpoints = []database.Endpoint{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpoints)
	}
}

func getEndpointHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

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

		endpoint, err := repo.GetEndpoint(r.Context(), id, userID)
		if err != nil {
			http.Error(w, `{"error":"endpoint not found"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(endpoint)
	}
}
