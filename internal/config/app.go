package config

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// AppContext holds shared infrastructure dependencies injected into every layer.
type AppContext struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

// Env reads an environment variable or returns a fallback default.
func Env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
