package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/Auto-Edge/autoedge-api/internal/config"
	"github.com/Auto-Edge/autoedge-api/internal/repository"
	"github.com/Auto-Edge/autoedge-api/internal/service"
	handler "github.com/Auto-Edge/autoedge-api/internal/transport/http"
)

func main() {
	// ---- Postgres ---------------------------------------------------------
	pgDSN := config.Env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/autoedge?sslmode=disable")

	pgPool, err := pgxpool.New(context.Background(), pgDSN)
	if err != nil {
		log.Fatalf("postgres: unable to create pool: %v", err)
	}
	defer pgPool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pgPool.Ping(ctx); err != nil {
		log.Fatalf("postgres: ping failed: %v", err)
	}
	log.Println("postgres: connected")

	// ---- Redis ------------------------------------------------------------
	redisAddr := config.Env("REDIS_URL", "localhost:6379")

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	defer rdb.Close()

	rctx, rcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer rcancel()
	if err := rdb.Ping(rctx).Err(); err != nil {
		log.Fatalf("redis: ping failed: %v", err)
	}
	log.Println("redis: connected")

	// ---- AWS S3 -----------------------------------------------------------
	awsCfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("aws: unable to load config: %v", err)
	}
	s3Client := s3.NewFromConfig(awsCfg)

	// ---- Wire layers (Repository -> Service -> Handler) -------------------
	modelRepo := repository.NewModelRepository(pgPool)
	storageSvc := service.NewStorageService(s3Client)

	s3Bucket := config.Env("S3_BUCKET", "autoedge-models")
	modelHandler := handler.NewModelHandler(modelRepo, storageSvc, rdb, s3Bucket)

	// ---- Fiber app --------------------------------------------------------
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	// Health check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Domain routes
	api := app.Group("/api/v1")
	handler.RegisterRoutes(api, modelHandler)

	// ---- Graceful shutdown ------------------------------------------------
	port := config.Env("PORT", "8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("fiber: listen error: %v", err)
		}
	}()
	log.Printf("server: listening on :%s", port)

	<-quit
	log.Println("server: shutting down …")

	if err := app.Shutdown(); err != nil {
		log.Fatalf("server: forced shutdown: %v", err)
	}
	log.Println("server: stopped")
}
