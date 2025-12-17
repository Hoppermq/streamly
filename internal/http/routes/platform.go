package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/internal/core/platform/organization"
	handlers "github.com/hoppermq/streamly/internal/http/handlers/platform"
)

func RegisterPlatformRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	organizationUseCase *organization.UseCase,
) {
	v1 := router.Group("/v1")
	organizationGroup := v1.Group("/organizations")
	{
		hndler := handlers.NewOrganization(
			handlers.OrganizationWithLogger(logger),
			handlers.OrganizationWithUseCase(organizationUseCase),
		)

		organizationGroup.GET("/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")
			logger.Info("performing get request", "id", id)

			ctx.JSON(http.StatusOK, gin.H{"id": id})
		})

		organizationGroup.POST("/", hndler.Create)
	}
}

func CreatePlatformRegistrar(
	logger *slog.Logger,
	orgUc *organization.UseCase,
) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterPlatformRoutes(engine, logger, orgUc)
	}
}
