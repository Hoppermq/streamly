package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hoppermq/streamly/internal/core/platform/organization"
	"github.com/hoppermq/streamly/pkg/domain"
)

type Organization struct {
	logger *slog.Logger
	uc     *organization.UseCase
}

type OrganizationOption func(*Organization) error

func OrganizationWithLogger(logger *slog.Logger) OrganizationOption {
	return func(o *Organization) error {
		o.logger = logger
		return nil
	}
}

func OrganizationWithUseCase(uc *organization.UseCase) OrganizationOption {
	return func(o *Organization) error {
		o.uc = uc
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

	// TODO: implement auth to extract jwt here.
	zitadelUserID := "354116995202968329"

	if err := o.uc.Create(c, org, zitadelUserID); err != nil {
		// TODO: should compare error from here.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Organization created successfully"})
}

func (o *Organization) Update(c *gin.Context) {
	var dataToUpdate domain.UpdateOrganization
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no parameter set"})
		return
	}

	if err := c.ShouldBind(&dataToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := o.uc.Update(c, id, dataToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization updated successfully"})
}

func (o *Organization) FindOneByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no parameter set"})
		return
	}

	org, err := o.uc.FindOneByID(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": org})
}

func (o *Organization) FindAll(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "100"
	}

	offset := c.Query("offset")
	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	orgs, err := o.uc.FindAll(c, limitInt, offsetInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orgs})
}

func (o *Organization) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no parameter set"})
		return
	}

	err := o.uc.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}
