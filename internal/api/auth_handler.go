package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
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

		req.Email = strings.TrimSpace(req.Email)
		if req.Email == "" || req.Password == "" {
			slog.Warn("failed signup attempt", "reason", "missing fields", "ip", r.RemoteAddr)
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusBadRequest)
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
			slog.Warn("failed signup attempt", "reason", "email already registered", "ip", r.RemoteAddr)
			// prevent enumeration: return same error as login
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
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

		slog.Info("successful signup", "email", req.Email, "ip", r.RemoteAddr)

		isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
		sameSite := http.SameSiteLaxMode
		if isSecure {
			sameSite = http.SameSiteNoneMode
		}
		
		http.SetCookie(w, &http.Cookie{
			Name:     "siphon_auth",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   isSecure,
			SameSite: sameSite,
			MaxAge:   24 * 60 * 60,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func loginHandler(repo *database.WebhookRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid JSON payload"}`, http.StatusBadRequest)
			return
		}

		req.Email = strings.TrimSpace(req.Email)
		if req.Email == "" || req.Password == "" {
			slog.Warn("failed login attempt", "reason", "missing fields", "ip", r.RemoteAddr)
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusBadRequest)
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
			slog.Warn("failed login attempt", "reason", "user not found", "ip", r.RemoteAddr)
			http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return
		}

		// verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			slog.Warn("failed login attempt", "reason", "incorrect password", "ip", r.RemoteAddr)
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

		slog.Info("successful login", "email", req.Email, "ip", r.RemoteAddr)

		isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
		sameSite := http.SameSiteLaxMode
		if isSecure {
			sameSite = http.SameSiteNoneMode
		}
		
		http.SetCookie(w, &http.Cookie{
			Name:     "siphon_auth",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   isSecure,
			SameSite: sameSite,
			MaxAge:   24 * 60 * 60,
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isSecure := r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
		sameSite := http.SameSiteLaxMode
		if isSecure {
			sameSite = http.SameSiteNoneMode
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "siphon_auth",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   isSecure,
			SameSite: sameSite,
			MaxAge:   -1,
		})
		slog.Info("successful logout", "ip", r.RemoteAddr)
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
