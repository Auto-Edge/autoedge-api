package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Auto-Edge/autoedge-api/internal/models"
)

// ModelRepository defines the contract. Any struct that implements these methods IS a ModelRepository.
type ModelRepository interface {
	CreateModelVersion(ctx context.Context, version *models.ModelVersion) error
}

// modelRepo is the concrete implementation holding the database connection
type modelRepo struct {
	db *pgxpool.Pool
}

// NewModelRepository is the constructor we will call in main.go
func NewModelRepository(db *pgxpool.Pool) ModelRepository {
	return &modelRepo{db: db}
}

// CreateModelVersion executes the raw SQL to save the pending job
func (r *modelRepo) CreateModelVersion(ctx context.Context, v *models.ModelVersion) error {
	query := `
		INSERT INTO model_versions (id, model_id, tag, status, s3_key, size_bytes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(ctx, query, v.ID, v.ModelID, v.Tag, v.Status, v.S3Key, v.SizeBytes, v.CreatedAt)
	return err
}