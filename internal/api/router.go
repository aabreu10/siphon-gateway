package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/aabreu10/siphon-gateway/internal/broker"
	"github.com/aabreu10/siphon-gateway/internal/config"
	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/aabreu10/siphon-gateway/internal/sse"
	"golang.org/x/time/rate"
)

// creates the router with middleware and routes
func NewRouter(repo *database.WebhookRepo, pub *broker.Publisher, hub *sse.Hub, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-Webhook-Source"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// health
	r.Get("/health", healthHandler())

	// api v1
	r.Route("/api/v1", func(r chi.Router) {
		// ingest route (protected by INGEST_API_KEY + Rate Limiter)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(cfg.IngestAPIKey))
			r.Use(RateLimitMiddleware(NewIPRateLimiter(rate.Limit(50), 100))) // 50 rps
			r.Post("/webhook", ingestHandler(repo, pub, hub, cfg.TargetURL))
		})
		
		// echo route
		r.Post("/echo", echoHandler())
		
		// dashboard admin routes (protected by ADMIN_API_KEY)
		r.Group(func(r chi.Router) {
			r.Use(AuthMiddleware(cfg.AdminAPIKey))
			r.Get("/webhooks", listHandler(repo))
			r.Get("/events", sseHandler(hub))
			r.Post("/webhook/{id}/replay", replayHandler(repo, pub, hub))
			r.Get("/webhook/{id}/logs", logsHandler(repo))
		})
	})

	return r
}
