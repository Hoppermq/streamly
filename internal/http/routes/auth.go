package routes

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/hoppermq/streamly/internal/http/handlers"
)

func RegisterAuthRoutes(
	router *gin.Engine,
	logger *slog.Logger,
) {
	authHandler := handlers.NewAuthHandler(
		handlers.AuthWithLogger(logger),
	)

	v1 := router.Group("/v1")
	authGroup := v1.Group("/auth")
	{
		userGroup := authGroup.Group("/user")
		{
			userGroup.POST("/login", authHandler.HandleUserLogin)
			userGroup.POST("/logout", authHandler.HandleUserLogout)
		}

		serviceGroup := authGroup.Group("/service")
		{
			serviceGroup.POST("/token", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"data": "token success"})
			})
		}
	}
}

func CreateAuthRegistrar(
	logger *slog.Logger,
) RouteRegistrar {
	return func(engine *gin.Engine) {
		RegisterAuthRoutes(engine, logger)
	}
}
