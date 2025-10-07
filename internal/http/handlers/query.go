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

func QueryHandlerWithLogger(logger *slog.Logger) QueryHandlerOption {
	return func(qh *QueryHandler) {
		qh.logger = logger
	}
}

func QueryHandlerWithUseCase(queryUseCase domain.QueryUseCase) QueryHandlerOption {
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

func (h *QueryHandler) PerformQuery(c *gin.Context) {
	res, err := h.useCase.AsyncQuery(c, nil)
	if err != nil {
		h.logger.Error("error performing query", "error", err.Error())
	}

	c.JSON(http.StatusOK, res)
}
