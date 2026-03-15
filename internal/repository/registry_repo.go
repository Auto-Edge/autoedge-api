package repository

import (
	"context"

	"github.com/Auto-Edge/autoedge-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegistryRepository defines the interface for model registry data access
type RegistryRepository interface {
	CreateModel(ctx context.Context, m *models.Model) error
	GetModelByID(ctx context.Context, id string) (*models.Model, error)
	ListModels(ctx context.Context, activeOnly bool) ([]models.Model, error)
	CreateModelVersion(ctx context.Context, v *models.ModelVersion) error
}

type PostgreRegistryRepo struct {
	db *pgxpool.Pool
}

func NewPostgreRegistryRepo(db *pgxpool.Pool) RegistryRepository {
	return &PostgreRegistryRepo{db: db}
}

func (r *PostgreRegistryRepo) CreateModel(ctx context.Context, m *models.Model) error {
	query := `
		INSERT INTO models (id, name, description, created_at, is_active)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, m.ID, m.Name, m.Description, m.CreatedAt, m.IsActive)
	return err
}

func (r *PostgreRegistryRepo) GetModelByID(ctx context.Context, id string) (*models.Model, error) {
	query := `SELECT id, name, description, created_at, is_active FROM models WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var m models.Model
	err := row.Scan(&m.ID, &m.Name, &m.Description, &m.CreatedAt, &m.IsActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *PostgreRegistryRepo) ListModels(ctx context.Context, activeOnly bool) ([]models.Model, error) {
	query := `SELECT id, name, description, created_at, is_active FROM models`
	if activeOnly {
		query += ` WHERE is_active = true`
	}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize as empty slice strictly to avoid nil JSON output
	modelsList := make([]models.Model, 0)
	for rows.Next() {
		var m models.Model
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.CreatedAt, &m.IsActive); err != nil {
			return nil, err
		}
		modelsList = append(modelsList, m)
	}
	return modelsList, nil
}

func (r *PostgreRegistryRepo) CreateModelVersion(ctx context.Context, v *models.ModelVersion) error {
	query := `
		INSERT INTO model_versions 
		(id, model_id, version, file_path, file_size_bytes, file_hash, precision, created_at, is_published, download_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		v.ID, v.ModelID, v.Version, v.FilePath, v.FileSizeBytes,
		v.FileHash, v.Precision, v.CreatedAt, v.IsPublished, v.DownloadCount,
	)
	return err
}
