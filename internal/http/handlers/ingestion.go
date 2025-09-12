package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/pkg/domain"
)

type IngestionHandler struct {
	ingestionUseCase domain.EventIngestionUseCase
}

func NewIngestionHandler(ingestionUseCase domain.EventIngestionUseCase) *IngestionHandler {
	return &IngestionHandler{
		ingestionUseCase: ingestionUseCase,
	}
}

func (h *IngestionHandler) IngestEvents(c *gin.Context) {
	var request domain.BatchIngestionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("IngestionHandler: JSON binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "validation_failed",
			"error":  "Invalid JSON payload: " + err.Error(),
		})
		return
	}

	log.Printf("IngestionHandler: Received ingestion request for tenant %s with %d events", 
		request.TenantID, len(request.Events))

	response, err := h.ingestionUseCase.IngestBatch(c.Request.Context(), &request)
	if err != nil {
		log.Printf("IngestionHandler: Use case failed: %v", err)
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, response)
	log.Printf("IngestionHandler: Successfully processed batch %s", response.BatchID)
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
		if containsKeyword(errorMsg, keyword) {
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
		if containsKeyword(errorMsg, keyword) {
			return true
		}
	}
	return false
}

func containsKeyword(text, keyword string) bool {
	return len(text) >= len(keyword) && 
		   text[:len(keyword)] == keyword || 
		   findSubstring(text, keyword)
}

func findSubstring(text, substring string) bool {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return true
		}
	}
	return false
}
