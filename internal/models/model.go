package models

import (
	"time"
)

// --- Database Entities ---

// Model matches the "models" table
type Model struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

// ModelVersion matches the "model_versions" table
type ModelVersion struct {
	ID            string    `json:"id" db:"id"`
	ModelID       string    `json:"model_id" db:"model_id"`
	Version       string    `json:"version" db:"version"`
	FilePath      string    `json:"file_path" db:"file_path"`
	FileSizeBytes int64     `json:"file_size_bytes" db:"file_size_bytes"`
	FileHash      string    `json:"file_hash" db:"file_hash"`
	Precision     string    `json:"precision" db:"precision"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	IsPublished   bool      `json:"is_published" db:"is_published"`
	DownloadCount int       `json:"download_count" db:"download_count"`
}

// --- API Request/Response DTOs ---

type CreateModelRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type CreateModelResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
}

type ModelListResponse struct {
	Models []Model `json:"models"`
	Total  int     `json:"total"`
}

type ConversionTaskPayload struct {
	ID        string `json:"id"`
	FilePath  string `json:"file_path"`
	IsDemo    bool   `json:"is_demo"`
	CreatedAt string `json:"created_at"`
}
