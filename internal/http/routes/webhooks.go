package routes

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	handlers "github.com/hoppermq/streamly/internal/http/handlers/platform"
)

func RegistrarWebhooksRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	userUseCase *user.UseCase, // keep it for later with a global webhook usecase.
) {
	v1 := router.Group("/v1")
	webhooks := v1.Group("/webhooks")

	userHandler, err := handlers.NewUser(
		handlers.UserWithLogger(logger),
		handlers.UserWithUC(userUseCase),
	)
	if err != nil {
		panic(err)
	}

	zitadelWebhooks := webhooks.Group("/zitadel")
	{
		zitadelWebhooks.POST("/user-created", userHandler.Create)
	}
}

func CreateWebhookRegistrar(
	logger *slog.Logger,
	userUseCase *user.UseCase,
) RouteRegistrar {
	return func(router *gin.Engine) {
		RegistrarWebhooksRoutes(router, logger, userUseCase)
	}
}
