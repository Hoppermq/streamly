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
	zitadelSignature := c.GetHeader("zitadel-signature")
	if zitadelSignature == "" {
		u.logger.Warn("invalid signature key")
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var payload domain.ZitadelEventUserCreated
	if err := c.ShouldBindJSON(&payload); err != nil {
		u.logger.Error("failed to parse webhook payload", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	u.logger.Info("webhook payload received", "payload", payload)
	if err := u.uc.CreateFromEvent(c, &payload); err != nil {
		u.logger.Warn("failed to create from event", "error", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created"})
}

func (u *User) Update(c *gin.Context) {}

func (u *User) Delete(c *gin.Context) {}
