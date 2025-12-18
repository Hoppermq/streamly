package routes

import (
	"log/slog"

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

		organizationGroup.POST("/", hndler.Create)

		organizationGroup.GET("/:id", hndler.FindOneByID)
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
