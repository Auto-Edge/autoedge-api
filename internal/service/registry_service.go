package service

import (
	"context"
	"time"

	"github.com/Auto-Edge/autoedge-api/internal/models"
	"github.com/Auto-Edge/autoedge-api/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

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
