package models

import "time"

type Conversion struct {
	ID           string    `json:"id" db:"id"`
	ModelID      *string   `json:"model_id,omitempty" db:"model_id"`
	Status       string    `json:"status" db:"status"` // pending, processing, completed, failed
	TargetFormat string    `json:"target_format" db:"target_format"`
	FilePath     string    `json:"file_path" db:"file_path"`
	IsDemo       bool      `json:"is_demo" db:"is_demo"`
	ResultPath   *string   `json:"result_path,omitempty" db:"result_path"`
	ErrorMessage *string   `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type StartConversionResponse struct {
	TaskID  string `json:"task_id"`
	Message string `json:"message"`
	S3Key   string `json:"s3_key"`
}

type ConversionListResponse struct {
	Conversions []Conversion `json:"conversions"`
	Total       int          `json:"total"`
}

type StatusResponse struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}
