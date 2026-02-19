package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoppermq/middles"
	middlesGin "github.com/hoppermq/middles/gin"
	"github.com/hoppermq/streamly/internal/core/platform/organization"
	handlers "github.com/hoppermq/streamly/internal/http/handlers/platform"
)

func RegisterPlatformRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	organizationUseCase *organization.UseCase,
	authMiddleware middles.Auth,
) {
	v1 := router.Group("/v1")
	v1.Use(middlesGin.RequireAuth(authMiddleware))
	organizationGroup := v1.Group("/organizations")
	{
		orgHandler := handlers.NewOrganization(
			handlers.OrganizationWithLogger(logger),
			handlers.OrganizationWithUseCase(organizationUseCase),
		)

		organizationGroup.GET("/", orgHandler.FindAll)
		organizationGroup.POST("/", orgHandler.Create)

		organizationGroup.GET("/:id", orgHandler.FindOneByID)
		organizationGroup.PATCH("/:id", orgHandler.Update)
		organizationGroup.DELETE("/:id", orgHandler.Delete)
	}

	membershipGroup := v1.Group("/memberships")
	{
		membershipGroup.POST("/add-user", func(context *gin.Context) {})
	}

	userGroup := v1.Group("/users")
	{
		userGroup.GET("/", func(context *gin.Context) {
			logger.Info("listing users")
			context.JSON(http.StatusOK, nil)
		})
		userGroup.POST("/new", func(context *gin.Context) {})

		userGroup.GET("/:id", func(context *gin.Context) {})
		userGroup.PATCH("/:id", func(context *gin.Context) {})
		userGroup.DELETE("/:id", func(context *gin.Context) {})
	}
}

func CreatePlatformRegistrar(
	logger *slog.Logger,
	orgUc *organization.UseCase,
	authMiddleware middles.Auth,
) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterPlatformRoutes(engine, logger, orgUc, authMiddleware)
	}
}
