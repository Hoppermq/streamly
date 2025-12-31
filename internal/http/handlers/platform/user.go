package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	"github.com/hoppermq/streamly/pkg/domain"
)

type User struct {
	logger *slog.Logger
	uc     *user.UseCase
}

type UserOption func(*User) error

func UserWithLogger(logger *slog.Logger) UserOption {
	return func(u *User) error {
		u.logger = logger
		return nil
	}
}

func UserWithUC(uc *user.UseCase) UserOption {
	return func(u *User) error {
		u.uc = uc
		return nil
	}
}

func NewUser(opts ...UserOption) (*User, error) {
	u := &User{}

	for _, opt := range opts {
		if err := opt(u); err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (u *User) FindOne(c *gin.Context) {
	var userInput domain.CreateUser
	if err := c.ShouldBind(&userInput); err != nil {
		err = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err := u.uc.Create(c, &userInput); err != nil {
		err = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (u *User) FindAll(c *gin.Context) {}

func (u *User) Create(c *gin.Context) {
	u.logger.Info("🎯 WEBHOOK RECEIVED - Zitadel user created")

	// Log the raw payload for debugging
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		u.logger.Error("failed to parse webhook payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	u.logger.Info("webhook payload received", "payload", payload)

	// TODO: Parse Zitadel payload and call u.uc.Create()

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (u *User) Update(c *gin.Context) {}

func (u *User) Delete(c *gin.Context) {}
