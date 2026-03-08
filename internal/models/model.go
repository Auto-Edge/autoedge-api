package models

import "time"

// Model represents the parent container (e.g., "FaceBlur")
type Model struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// ModelVersion represents a specific build (e.g., "v1.5.0")
type ModelVersion struct {
	ID         string    `json:"id"`
	ModelID    string    `json:"model_id"`
	Tag        string    `json:"tag"`
	Status     string    `json:"status"` // "pending", "ready", "failed"
	S3Key      string    `json:"s3_key"`
	SchemaHash string    `json:"schema_hash,omitempty"` // Extracted later by Python
	SizeBytes  int64     `json:"size_bytes"`
	CreatedAt  time.Time `json:"created_at"`
}