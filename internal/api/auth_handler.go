package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

var jwtSecret = []byte("super-secret-jwt-key") // In production, use env var

func generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func signupHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, `{"error":"email and password are required"}`, http.StatusBadRequest)
			return
		}

		// check if user exists
		existing, err := repo.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			slog.Error("error checking existing user", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if existing != nil {
			http.Error(w, `{"error":"email already registered"}`, http.StatusConflict)
			return
		}

		// hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error("error hashing password", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// create user
		userID, err := repo.CreateUser(r.Context(), req.Email, string(hash))
		if err != nil {
			slog.Error("error creating user", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		// generate token
		token, err := generateJWT(userID.String())
		if err != nil {
			slog.Error("error generating jwt", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(AuthResponse{Token: token})
	}
}

func loginHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, `{"error":"email and password are required"}`, http.StatusBadRequest)
			return
		}

		// get user
		user, err := repo.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			slog.Error("error finding user", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// generate token
		token, err := generateJWT(user.ID.String())
		if err != nil {
			slog.Error("error generating jwt", "error", err)
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthResponse{Token: token})
	}
}
