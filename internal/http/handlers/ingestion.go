package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/internal/common"
	"github.com/hoppermq/streamly/pkg/domain"
)

type IngestionHandler struct {
	ingestionUseCase domain.EventIngestionUseCase

	logger *slog.Logger
}

type Option func(*IngestionHandler)

func WithLogger(logger *slog.Logger) Option {
	return func(h *IngestionHandler) {
		h.logger = logger
	}
}

func WithUSeCase(ingestionUseCase domain.EventIngestionUseCase) Option {
	return func(h *IngestionHandler) {
		h.ingestionUseCase = ingestionUseCase
	}
}

func NewIngestionHandler(opts ...Option) *IngestionHandler {
	h := &IngestionHandler{}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *IngestionHandler) IngestEvents(c *gin.Context) {
	var request domain.BatchIngestionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("json ingestion failed", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	h.logger.Info("ingestion request", "request", request)

	response, err := h.ingestionUseCase.IngestBatch(c.Request.Context(), &request)
	if err != nil {
		h.logger.Warn("ingestion failed", "error", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, response)
	h.logger.Info("ingestion succeeded", "response", response)
}

func (h *IngestionHandler) handleError(c *gin.Context, err error) {
	errorMsg := err.Error()

	if isValidationError(errorMsg) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "validation_failed",
			"error":  errorMsg,
		})
		return
	}

	if isRepositoryError(errorMsg) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "service_unavailable",
			"error":  "Storage service temporarily unavailable",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "internal_error",
		"error":  "An unexpected error occurred",
	})
}

func isValidationError(errorMsg string) bool {
	validationKeywords := []string{
		"validation failed",
		"is required",
		"cannot be empty",
		"cannot exceed",
		"must be valid JSON",
	}

	for _, keyword := range validationKeywords {
		if common.ContainsKeyword(errorMsg, keyword) {
			return true
		}
	}
	return false
}

func isRepositoryError(errorMsg string) bool {
	repositoryKeywords := []string{
		"repository insert failed",
		"connection refused",
		"timeout",
		"database",
	}

	for _, keyword := range repositoryKeywords {
		if common.ContainsKeyword(errorMsg, keyword) {
			return true
		}
	}
	return false
}
