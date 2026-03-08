package http

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Auto-Edge/autoedge-api/internal/models"
	"github.com/Auto-Edge/autoedge-api/internal/repository"
	"github.com/Auto-Edge/autoedge-api/internal/service"
)

// ModelHandler is the transport layer for model-related endpoints.
type ModelHandler struct {
	repo    repository.ModelRepository
	storage service.StorageService
	redis   *redis.Client
	bucket  string
}

// NewModelHandler constructs a ModelHandler with all required dependencies.
func NewModelHandler(repo repository.ModelRepository, storage service.StorageService, rdb *redis.Client, bucket string) *ModelHandler {
	return &ModelHandler{
		repo:    repo,
		storage: storage,
		redis:   rdb,
		bucket:  bucket,
	}
}

// RegisterRoutes mounts the model endpoints onto a Fiber router group.
func RegisterRoutes(r fiber.Router, h *ModelHandler) {
	r.Post("/models/upload-url", h.CreateUploadURL)
}

// --- Request / Response DTOs ------------------------------------------------

type createUploadURLRequest struct {
	ModelID   string `json:"model_id"`
	Tag       string `json:"tag"`
	Filename  string `json:"filename"`
	SizeBytes int64  `json:"size_bytes"`
}

type createUploadURLResponse struct {
	VersionID string `json:"version_id"`
	UploadURL string `json:"upload_url"`
	S3Key     string `json:"s3_key"`
}

// --- Handler ----------------------------------------------------------------

const redisQueueKey = "queue:model_conversion"

// CreateUploadURL handles POST /api/v1/models/upload-url.
// It generates a presigned S3 PUT URL, persists a pending ModelVersion,
// and enqueues a conversion job for the Python Celery worker.
func (h *ModelHandler) CreateUploadURL(c *fiber.Ctx) error {
	var req createUploadURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid JSON body",
		})
	}

	// ---- Validate required fields -----------------------------------------
	if req.ModelID == "" || req.Tag == "" || req.Filename == "" || req.SizeBytes <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "model_id, tag, filename, and a positive size_bytes are required",
		})
	}

	versionID := uuid.New().String()
	s3Key := fmt.Sprintf("models/%s/%s/%s", req.ModelID, versionID, req.Filename)

	// ---- 1. Generate S3 Pre-signed URL ------------------------------------
	uploadURL, err := h.storage.GeneratePresignedUploadURL(c.Context(), h.bucket, s3Key, 15)
	if err != nil {
		log.Printf("error: s3 presign failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate upload URL",
		})
	}

	// ---- 2. Persist pending ModelVersion -----------------------------------
	version := &models.ModelVersion{
		ID:        versionID,
		ModelID:   req.ModelID,
		Tag:       req.Tag,
		Status:    "pending",
		S3Key:     s3Key,
		SizeBytes: req.SizeBytes,
		CreatedAt: time.Now().UTC(),
	}	

	if err := h.repo.CreateModelVersion(c.Context(), version); err != nil {
		log.Printf("error: db insert failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create model version",
		})
	}

	// ---- 3. Enqueue conversion job to Redis for Celery --------------------
	job, _ := json.Marshal(fiber.Map{
		"version_id": versionID,
		"s3_key":     s3Key,
	})

	if err := h.redis.RPush(c.Context(), redisQueueKey, job).Err(); err != nil {
		log.Printf("error: redis rpush failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to enqueue conversion job",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createUploadURLResponse{
		VersionID: versionID,
		UploadURL: uploadURL,
		S3Key:     s3Key,
	})
}
