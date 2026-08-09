package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	// Use standard localhost:5432 which is mapped in docker-compose
	pool, err := pgxpool.New(ctx, "postgres://siphon:siphon_secret@localhost:5432/siphon_gateway?sslmode=disable")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM webhooks").Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	fmt.Printf("Total webhooks in DB: %d\n", count)

	rows, err := pool.Query(ctx, "SELECT id, source, endpoint_id, target_url FROM webhooks LIMIT 10")
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, source, target_url string
		var endpoint_id interface{}
		err = rows.Scan(&id, &source, &endpoint_id, &target_url)
		if err != nil {
			log.Fatalf("Scan failed: %v\n", err)
		}
		fmt.Printf("Webhook ID: %s, Source: %s, Endpoint: %v, Target: %s\n", id, source, endpoint_id, target_url)
	}
}
