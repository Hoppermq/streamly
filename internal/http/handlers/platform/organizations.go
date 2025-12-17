package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hoppermq/streamly/pkg/domain"
)

type Organization struct {
	logger *slog.Logger
	uc     any
}

type OrganizationOption func(*Organization) error

func OrganizationWithLogger(logger *slog.Logger) OrganizationOption {
	return func(o *Organization) error {
		o.logger = logger
		return nil
	}
}

func NewOrganization(opts ...OrganizationOption) *Organization {
	o := &Organization{}

	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil
		}
	}

	return o
}

func (o *Organization) Create(c *gin.Context) {
	var org domain.CreateOrganization
	if err := c.ShouldBind(&org); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Organization created successfully"})
}
