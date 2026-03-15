package http

import (
	"github.com/Auto-Edge/autoedge-api/internal/models"
	"github.com/Auto-Edge/autoedge-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type RegistryHandler struct {
	svc     *service.RegistryService
	storage *service.StorageService
}

func NewRegistryHandler(svc *service.RegistryService, storage *service.StorageService) *RegistryHandler {
	return &RegistryHandler{
		svc:     svc,
		storage: storage,
	}
}

// RegisterRoutes maps the endpoints
func (h *RegistryHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	registry := api.Group("/registry")
	registry.Post("/models", h.CreateModel)
	registry.Get("/models", h.ListModels)
	registry.Get("/models/:id", h.GetModel)

	conversion := api.Group("/conversion")
	conversion.Post("/upload", h.UploadModel)
}

func (h *RegistryHandler) CreateModel(c *fiber.Ctx) error {
	var req models.CreateModelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Basic Validation (in a real app, use GoValidator)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	model, err := h.svc.CreateModel(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(model)
}

func (h *RegistryHandler) ListModels(c *fiber.Ctx) error {
	// Simple mapping of query param
	activeOnly := c.Query("active_only", "true") == "true"

	// Delegate to service
	// Note: We use c.Context() or c.UserContext() for timeouts
	result, err := h.svc.ListModels(c.Context(), activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch models"})
	}

	return c.JSON(fiber.Map{"models": result, "total": len(result)})
}

func (h *RegistryHandler) GetModel(c *fiber.Ctx) error {
	id := c.Params("id")
	model, err := h.svc.GetModel(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if model == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Model not found"})
	}
	return c.JSON(model)
}

func (h *RegistryHandler) UploadModel(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File upload required"})
	}

	// Open the multipart stream
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file stream"})
	}
	defer file.Close()

	// Upload to S3 via Storage Service
	key, err := h.storage.UploadFile(c.Context(), file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to storage"})
	}

	// Trigger the "Background Task" logic via Service
	taskID, err := h.svc.StartConversion(c.Context(), key, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to queue conversion"})
	}

	return c.JSON(fiber.Map{
		"task_id": taskID,
		"message": "File uploaded and conversion queued",
		"s3_key":  key,
	})
}
