package database

import (
	"context"
	"testing"

	"github.com/ericlamnguyen/distributed-job-queue/internal/config"
)

func TestNewPool(t *testing.T) {
	cfg := config.Load()
	databaseURL := cfg.DatabaseURL

	ctx := context.Background()

	pool, err := NewPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()
}
