package api

import (
	"context"
	"net/http"
	"strings"
	"sync"
	
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/time/rate"
)

// rate limiter per IP
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(limiter *IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract IP
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = strings.Split(forwarded, ",")[0]
			}

			if !limiter.GetLimiter(ip).Allow() {
				http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type contextKey string

const UserIDKey contextKey = "user_id"

// JWT Auth middleware
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var providedToken string
			
			// 1. Check cookie first
			if cookie, err := r.Cookie("siphon_auth"); err == nil {
				providedToken = cookie.Value
			}

			// 2. Fallback to Authorization header
			authHeader := r.Header.Get("Authorization")
			if providedToken == "" && authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				providedToken = strings.TrimPrefix(authHeader, "Bearer ")
			}

			// 3. Fallback to query param (for SSE)
			if providedToken == "" && r.URL.Path == "/api/v1/events" {
				providedToken = r.URL.Query().Get("api_key")
			}

			if providedToken == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
				if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil // using jwtSecret from auth_handler.go
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			if _, ok := claims["exp"]; !ok {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			userIDStr, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userIDStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IngestAuthMiddleware protects the ingest route with a static key
func IngestAuthMiddleware(expectedKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if expectedKey == "" {
				next.ServeHTTP(w, r)
				return
			}

			// 1. Try static API key
			authHeader := r.Header.Get("Authorization")
			var providedKey string
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				providedKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				providedKey = r.URL.Query().Get("api_key")
			}

			if providedKey == expectedKey {
				next.ServeHTTP(w, r)
				return
			}

			// 2. Try JWT Cookie (allows dashboard simulator to hit this endpoint)
			if cookie, err := r.Cookie("siphon_auth"); err == nil && cookie.Value != "" {
				token, _ := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
					return jwtSecret, nil
				})
				if token != nil && token.Valid {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		})
	}
}

// HSTSMiddleware adds Strict-Transport-Security header
func HSTSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			next.ServeHTTP(w, r)
		})
	}
}

// MaxBytesMiddleware limits the size of incoming requests
func MaxBytesMiddleware(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

// CSRFMiddleware strictly enforces Content-Type: application/json for mutating requests
func CSRFMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				contentType := r.Header.Get("Content-Type")
				if !strings.HasPrefix(contentType, "application/json") {
					http.Error(w, `{"error":"invalid content type, must be application/json"}`, http.StatusUnsupportedMediaType)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
