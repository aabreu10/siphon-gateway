package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aabreu10/siphon-gateway/internal/api"
	"github.com/aabreu10/siphon-gateway/internal/broker"
	"github.com/aabreu10/siphon-gateway/internal/config"
	"github.com/aabreu10/siphon-gateway/internal/database"
	"github.com/aabreu10/siphon-gateway/internal/sse"
	"github.com/aabreu10/siphon-gateway/internal/worker"
)

func main() {
	// json logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting Siphon Gateway")

	// load config
	cfg := config.Load()

	// root context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ── postgresql ──────────────────────────────────────────────────────
	pool, err := database.NewPool(ctx, cfg.PostgresDSN())
	if err != nil {
		slog.Error("failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := database.NewWebhookRepo(pool)

	// ── rabbitmq ────────────────────────────────────────────────────────
	rmqConn, err := broker.Connect(cfg.RabbitMQURL)
	if err != nil {
		slog.Error("failed to connect to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer rmqConn.Close()

	// channel for topology and publishing
	pubCh, err := rmqConn.Channel()
	if err != nil {
		slog.Error("failed to open RabbitMQ channel", "error", err)
		os.Exit(1)
	}
	defer pubCh.Close()

	if err := broker.DeclareTopology(pubCh); err != nil {
		slog.Error("failed to declare RabbitMQ topology", "error", err)
		os.Exit(1)
	}

	pub := broker.NewPublisher(pubCh)

	// ── sse hub ─────────────────────────────────────────────────────────
	hub := sse.NewHub()

	// ── worker pool ─────────────────────────────────────────────────────
	if err := worker.StartWorkers(ctx, cfg.WorkerConcurrency, rmqConn, repo, pub, hub); err != nil {
		slog.Error("failed to start worker pool", "error", err)
		os.Exit(1)
	}

	// ── http server ─────────────────────────────────────────────────────
	router := api.NewRouter(repo, pub, hub, cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0, // disabled for sse streaming
		IdleTimeout:  60 * time.Second,
	}

	// start server in background
	go func() {
		slog.Info("HTTP server listening", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// ── graceful shutdown ───────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutdown signal received, draining...")

	// cancel workers
	cancel()

	// 10s drain timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server forced to shutdown", "error", err)
	}

	slog.Info("Siphon Gateway stopped")
}
