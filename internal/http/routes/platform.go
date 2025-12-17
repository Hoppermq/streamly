package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPlatformRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	uc ...any,
) {
	v1 := router.Group("/v1")
	organizationGroup := v1.Group("/organizations")
	{
		organizationGroup.GET("/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			logger.Info("performing get request", "id", id)

			ctx.JSON(http.StatusOK, gin.H{"id": id})
		})
	}
}

func CreatePlatformRegistrar(
	logger *slog.Logger,
) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterPlatformRoutes(engine, logger)
	}
}
