package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Auto-Edge/autoedge-api/internal/models"
	"github.com/Auto-Edge/autoedge-api/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const ConversionQueue = "conversion_tasks"

type ConversionService struct {
	repo  repository.ConversionRepository
	redis *redis.Client
}

func NewConversionService(repo repository.ConversionRepository, rdb *redis.Client) *ConversionService {
	return &ConversionService{
		repo:  repo,
		redis: rdb,
	}
}

func (s *ConversionService) StartConversion(ctx context.Context, filePath string, isDemo bool) (string, error) {
	taskID := uuid.New().String()

	// 1. Create DB Record
	conversion := &models.Conversion{
		ID:           taskID,
		Status:       "pending",
		TargetFormat: "onnx", // Default or explicit param
		FilePath:     filePath,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateConversion(ctx, conversion); err != nil {
		return "", fmt.Errorf("failed to create conversion record: %w", err)
	}

	// 2. Create Payload
	payload := models.ConversionTaskPayload{
		ID:        taskID,
		FilePath:  filePath,
		IsDemo:    isDemo,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// 3. Push to Redis
	err = s.redis.RPush(ctx, ConversionQueue, jsonPayload).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue task: %w", err)
	}

	return taskID, nil
}

func (s *ConversionService) GetConversionStatus(ctx context.Context, id string) (*models.Conversion, error) {
	return s.repo.GetConversion(ctx, id)
}
