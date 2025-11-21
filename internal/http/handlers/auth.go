package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	logger *slog.Logger
}

type AuthOption func(*AuthHandler)

func AuthWithLogger(logger *slog.Logger) AuthOption {
	return func(h *AuthHandler) {
		h.logger = logger
	}
}

func NewAuthHandler(opts ...AuthOption) *AuthHandler {
	h := &AuthHandler{}
	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *AuthHandler) HandleUserLogin(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *AuthHandler) HandleUserLogout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
