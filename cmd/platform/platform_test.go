//go:build integration

package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/core/platform/organization"
	handlers "github.com/hoppermq/streamly/internal/http/handlers/platform"
	"github.com/hoppermq/streamly/internal/models"
	"github.com/hoppermq/streamly/internal/storage/postgres"
	"github.com/hoppermq/streamly/internal/tests/testcontainers"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(t *testing.T, pgContainer *testcontainers.PostgresContainer) *gin.Engine {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	pgClient := postgres.NewClient(
		postgres.WithDB(pgContainer.BunDB),
		postgres.WithLogger(logger),
	)
	err := pgClient.Bootstrap(context.Background())
	require.NoError(t, err)

	repo, err := organization.NewRepository(
		organization.RepositoryWithDB(pgContainer.BunDB),
		organization.RepositoryWithLogger(logger),
	)
	require.NoError(t, err)

	uc, err := organization.NewUseCase(
		organization.UseCaseWithRepository(repo),
		organization.UseCaseWithLogger(logger),
		organization.UseCaseWithGenerator(uuid.New),
	)
	require.NoError(t, err)

	handler := handlers.NewOrganization(
		handlers.OrganizationWithUseCase(uc),
		handlers.OrganizationWithLogger(logger),
	)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/api/v1/organizations", handler.Create)
	router.GET("/api/v1/organizations", handler.FindAll)
	router.GET("/api/v1/organizations/:id", handler.FindOneByID)
	router.PUT("/api/v1/organizations/:id", handler.Update)
	router.DELETE("/api/v1/organizations/:id", handler.Delete)

	return router
}

func TestPlatformService_OrganizationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	pgContainer, err := testcontainers.StartPostgres(ctx)
	require.NoError(t, err)
	defer pgContainer.Close(ctx)

	router := setupTestRouter(t, pgContainer)

	t.Run("Create Organization via API", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":     "Test Org",
			"metadata": map[string]string{"key": "value"},
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/organizations", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var orgs []models.Organization
		err := pgContainer.BunDB.NewSelect().
			Model(&orgs).
			Where("name = ?", "Test Org").
			Scan(ctx)
		require.NoError(t, err)
		assert.Len(t, orgs, 1)
		assert.Equal(t, "Test Org", orgs[0].Name)
	})

	t.Run("List Organizations via API", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/organizations", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &result)
		require.NoError(t, err)

		data, ok := result["data"].([]interface{})
		assert.True(t, ok)
		assert.GreaterOrEqual(t, len(data), 1)
	})

	t.Run("Get Organization by ID via API", func(t *testing.T) {
		domainOrg := testutil.NewSampleOrganization("GetByID Test Org")
		modelOrg := models.Organization{
			Identifier: domainOrg.Identifier,
			Name:       domainOrg.Name,
			CreatedAt:  domainOrg.CreatedAt,
			UpdatedAt:  domainOrg.UpdatedAt,
		}
		_, err := pgContainer.BunDB.NewInsert().Model(&modelOrg).Exec(ctx)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/organizations/"+modelOrg.Identifier.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result map[string]interface{}
		err = json.Unmarshal(resp.Body.Bytes(), &result)
		require.NoError(t, err)

		data := result["data"].(map[string]interface{})
		assert.Equal(t, modelOrg.Name, data["Name"])
	})

	t.Run("Update Organization via API", func(t *testing.T) {
		domainOrg := testutil.NewSampleOrganization("Update Test Org")
		modelOrg := models.Organization{
			Identifier: domainOrg.Identifier,
			Name:       domainOrg.Name,
			CreatedAt:  domainOrg.CreatedAt,
			UpdatedAt:  domainOrg.UpdatedAt,
		}
		_, err := pgContainer.BunDB.NewInsert().Model(&modelOrg).Exec(ctx)
		require.NoError(t, err)

		payload := map[string]interface{}{
			"name": "Updated Org Name",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/organizations/"+modelOrg.Identifier.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Logf("Response body: %s", resp.Body.String())
		}
		assert.Equal(t, http.StatusOK, resp.Code)

		var updatedOrg models.Organization
		err = pgContainer.BunDB.NewSelect().
			Model(&updatedOrg).
			Where("identifier = ?", modelOrg.Identifier).
			Scan(ctx)
		require.NoError(t, err)
		assert.Equal(t, "Updated Org Name", updatedOrg.Name)
	})

	t.Run("Delete Organization via API", func(t *testing.T) {
		domainOrg := testutil.NewSampleOrganization("Delete Test Org")
		modelOrg := models.Organization{
			Identifier: domainOrg.Identifier,
			Name:       domainOrg.Name,
			CreatedAt:  domainOrg.CreatedAt,
			UpdatedAt:  domainOrg.UpdatedAt,
		}
		_, err := pgContainer.BunDB.NewInsert().Model(&modelOrg).Exec(ctx)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/organizations/"+modelOrg.Identifier.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}

func TestPlatformService_DatabaseIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	pgContainer, err := testcontainers.StartPostgres(ctx)
	require.NoError(t, err)
	defer pgContainer.Close(ctx)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	pgClient := postgres.NewClient(
		postgres.WithDB(pgContainer.BunDB),
		postgres.WithLogger(logger),
	)
	err = pgClient.Bootstrap(ctx)
	require.NoError(t, err)

	t.Run("Repository Insert and Select", func(t *testing.T) {
		repo, err := organization.NewRepository(
			organization.RepositoryWithDB(pgContainer.BunDB),
			organization.RepositoryWithLogger(logger),
		)
		require.NoError(t, err)

		org := &domain.Organization{
			Identifier: uuid.New(),
			Name:       "Repository Test Org",
		}

		err = repo.Create(ctx, org)
		require.NoError(t, err)

		result, err := repo.FindOneByID(ctx, org.Identifier)
		require.NoError(t, err)
		assert.Equal(t, org.Name, result.Name)
		assert.Equal(t, org.Identifier, result.Identifier)
	})

	t.Run("Repository FindAll with Pagination", func(t *testing.T) {
		repo, err := organization.NewRepository(
			organization.RepositoryWithDB(pgContainer.BunDB),
			organization.RepositoryWithLogger(logger),
		)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			org := &domain.Organization{
				Identifier: uuid.New(),
				Name:       "Pagination Test Org " + uuid.NewString(),
			}
			err = repo.Create(ctx, org)
			require.NoError(t, err)
		}

		orgs, err := repo.FindAll(ctx, 3, 0)
		require.NoError(t, err)
		assert.Len(t, orgs, 3)

		orgsPage2, err := repo.FindAll(ctx, 3, 3)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(orgsPage2), 2)
	})
}
