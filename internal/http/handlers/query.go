package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/pkg/domain"
)

type QueryHandler struct {
	useCase domain.QueryUseCase
	logger  *slog.Logger
}

type QueryHandlerOption func(*QueryHandler)

func WithQueryLogger(logger *slog.Logger) QueryHandlerOption {
	return func(qh *QueryHandler) {
		qh.logger = logger
	}
}

func WithQueryUseCase(queryUseCase domain.QueryUseCase) QueryHandlerOption {
	return func(qh *QueryHandler) {
		qh.useCase = queryUseCase
	}
}

func NewQueryHandler(options ...QueryHandlerOption) *QueryHandler {
	qh := &QueryHandler{}
	for _, opt := range options {
		opt(qh)
	}
	return qh
}

func (h *QueryHandler) Execute(c *gin.Context) {
	var req domain.QueryAstRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Info("error while binding query request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.TenantID = extractTenantID(c)

	res, err := h.useCase.SyncQuery(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("error performing query", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, res)
}


func extractTenantID(ctx *gin.Context) string {
	tenantID, exists := ctx.Get("tenant_id")
	if !exists {
		return "default"
	}
	if tenantIDStr, ok := tenantID.(string); ok {
		return tenantIDStr
	}
	return "default"
}

