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

type RegistryService struct {
	repo  repository.RegistryRepository
	redis *redis.Client
}

func NewRegistryService(repo repository.RegistryRepository, rdb *redis.Client) *RegistryService {
	return &RegistryService{
		repo:  repo,
		redis: rdb,
	}
}

// CreateModel implements the business logic for creating a model
func (s *RegistryService) CreateModel(ctx context.Context, req models.CreateModelRequest) (*models.Model, error) {
	model := &models.Model{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		IsActive:    true,
	}

	if err := s.repo.CreateModel(ctx, model); err != nil {
		return nil, err
	}

	return model, nil
}

// ListModels returns a list of models
func (s *RegistryService) ListModels(ctx context.Context, activeOnly bool) ([]models.Model, error) {
	return s.repo.ListModels(ctx, activeOnly)
}

// GetModel returns a model by ID
func (s *RegistryService) GetModel(ctx context.Context, id string) (*models.Model, error) {
	return s.repo.GetModelByID(ctx, id)
}

// StartConversion handles the worker handoff replacing FastAPI BackgroundTasks
func (s *RegistryService) StartConversion(ctx context.Context, filePath string, isDemo bool) (string, error) {
	taskID := uuid.New().String()

	// Create the payload for the Python worker
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

	// Push to Redis Queue instead of calling Celery directly
	// The Python worker must be configured to pop from this list
	err = s.redis.RPush(ctx, ConversionQueue, jsonPayload).Err()
	if err != nil {
		return "", fmt.Errorf("failed to enqueue task: %w", err)
	}

	return taskID, nil
}
