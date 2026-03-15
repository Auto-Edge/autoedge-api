package repository

import (
	"context"

	"github.com/Auto-Edge/autoedge-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConversionRepository interface {
	CreateConversion(ctx context.Context, c *models.Conversion) error
	GetConversion(ctx context.Context, id string) (*models.Conversion, error)
}

type conversionRepo struct {
	db *pgxpool.Pool
}

func NewConversionRepo(db *pgxpool.Pool) ConversionRepository {
	return &conversionRepo{db: db}
}

func (r *conversionRepo) CreateConversion(ctx context.Context, c *models.Conversion) error {
	query := `
        INSERT INTO conversions (id, model_version_id, status, target_format, input_path, output_path, error_message, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := r.db.Exec(ctx, query,
		c.ID, c.ModelID, c.Status, c.TargetFormat, c.FilePath, c.IsDemo, c.ResultPath, c.ErrorMessage, c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *conversionRepo) GetConversion(ctx context.Context, id string) (*models.Conversion, error) {
	query := `
        SELECT id, model_version_id, status, target_format, input_path, output_path, error_message, created_at, updated_at, finished_at 
        FROM conversions WHERE id = $1
    `
	row := r.db.QueryRow(ctx, query, id)

	var c models.Conversion
	err := row.Scan(
		&c.ID, &c.ModelID, &c.Status, &c.TargetFormat, &c.FilePath, &c.IsDemo, &c.ResultPath,
		&c.ErrorMessage, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}
