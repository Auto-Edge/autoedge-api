package http

import (
	"github.com/Auto-Edge/autoedge-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ConversionHandler struct {
	svc     *service.ConversionService
	storage *service.StorageService
}

func NewConversionHandler(svc *service.ConversionService, storage *service.StorageService) *ConversionHandler {
	return &ConversionHandler{
		svc:     svc,
		storage: storage,
	}
}

func (h *ConversionHandler) RegisterRoutes(api fiber.Router) {
	conversion := api.Group("/conversions")
	conversion.Post("/upload", h.UploadModel)
	conversion.Get("/:id", h.GetConversionStatus)
}

func (h *ConversionHandler) UploadModel(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File upload required"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file stream"})
	}
	defer file.Close()

	key, err := h.storage.UploadFile(c.Context(), file, fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload file to storage"})
	}

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

func (h *ConversionHandler) GetConversionStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Conversion ID is required"})
	}

	conversion, err := h.svc.GetConversionStatus(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch status"})
	}
	if conversion == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Conversion not found"})
	}

	return c.JSON(conversion)
}
