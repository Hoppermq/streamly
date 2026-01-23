package organization_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/core/platform/organization"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCaseCreate(t *testing.T) {
	t.Parallel()

	fixedUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	userUUID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name      string
		input     domain.CreateOrganization
		setupMock func(*mocks.MockUnitOfWork, *mocks.MockOrganizationRepository, *mocks.MockUserRepository, *mocks.MockMembershipRepository)
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name: "success - creates organization with generated UUID",
			input: domain.CreateOrganization{
				Name:     "Acme Corp",
				Metadata: map[string]string{"industry": "tech"},
			},
			setupMock: func(uow *mocks.MockUnitOfWork, orgRepo *mocks.MockOrganizationRepository, userRepo *mocks.MockUserRepository, memRepo *mocks.MockMembershipRepository) {
				uow.EXPECT().Organization().Return(orgRepo).Times(1)
				uow.EXPECT().User().Return(userRepo).Times(1)
				uow.EXPECT().Membership().Return(memRepo).Times(1)
				uow.EXPECT().Commit().Return(nil).Once()

				orgRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(org *domain.Organization) bool {
						return org.Name == "Acme Corp" && org.Identifier == fixedUUID
					})).
					Return(nil).
					Once()

				userRepo.EXPECT().
					GetUserIDFromZitadelID(mock.Anything, "test-zitadel-user-id").
					Return(userUUID, nil).
					Once()

				memRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(m *domain.Membership) bool {
						return m.OrgIdentifier == fixedUUID && m.UserIdentifier == userUUID
					})).
					Return(nil).
					Once()
			},
			assertErr: assert.NoError,
		},
		{
			name: "error - repository fails to insert",
			input: domain.CreateOrganization{
				Name:     "Test Org",
				Metadata: map[string]string{"key": "value"},
			},
			setupMock: func(uow *mocks.MockUnitOfWork, orgRepo *mocks.MockOrganizationRepository, userRepo *mocks.MockUserRepository, memRepo *mocks.MockMembershipRepository) {
				uow.EXPECT().Organization().Return(orgRepo).Times(1)
				uow.EXPECT().Rollback().Return(nil).Once()

				orgRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(errors.New("database connection failed")).
					Once()
			},
			assertErr: assert.Error,
		},
		{
			name: "success - handles empty metadata map",
			input: domain.CreateOrganization{
				Name:     "Minimal Org",
				Metadata: map[string]string{},
			},
			setupMock: func(uow *mocks.MockUnitOfWork, orgRepo *mocks.MockOrganizationRepository, userRepo *mocks.MockUserRepository, memRepo *mocks.MockMembershipRepository) {
				uow.EXPECT().Organization().Return(orgRepo).Times(1)
				uow.EXPECT().User().Return(userRepo).Times(1)
				uow.EXPECT().Membership().Return(memRepo).Times(1)
				uow.EXPECT().Commit().Return(nil).Once()

				orgRepo.EXPECT().
					Create(mock.Anything, mock.MatchedBy(func(org *domain.Organization) bool {
						return org.Name == "Minimal Org" && org.Identifier == fixedUUID
					})).
					Return(nil).
					Once()

				userRepo.EXPECT().
					GetUserIDFromZitadelID(mock.Anything, "test-zitadel-user-id").
					Return(userUUID, nil).
					Once()

				memRepo.EXPECT().
					Create(mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			assertErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uow := mocks.NewMockUnitOfWork(t)
			uowFactory := mocks.NewMockUnitOfWorkFactory(t)
			orgRepo := mocks.NewMockOrganizationRepository(t)
			userRepo := mocks.NewMockUserRepository(t)
			memRepo := mocks.NewMockMembershipRepository(t)

			uowFactory.EXPECT().NewUnitOfWork(mock.Anything).Return(uow, nil).Once()
			tt.setupMock(uow, orgRepo, userRepo, memRepo)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			ctx := context.Background()

			uc, err := organization.NewUseCase(
				organization.UseCaseWithGenerator(func() uuid.UUID { return fixedUUID }),
				organization.UseCaseWithLogger(logger),
				organization.UseCaseWithRepository(orgRepo),
				organization.UseCaseWithUOW(uowFactory),
				organization.UseCaseWithUUIDParser(uuid.Parse),
			)
			require.NoError(t, err)

			err = uc.Create(ctx, tt.input, "test-zitadel-user-id")

			tt.assertErr(t, err)
		})
	}
}

func TestUseCaseFindOneByID(t *testing.T) {
	t.Parallel()

	orgID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	nonExistingOrgID := uuid.MustParse("123e4567-e89b-12d3-a456-42661417400e")
	expectedOrg := &domain.Organization{
		Identifier: orgID,
		Name:       "Test Org",
	}

	tests := []struct {
		name        string
		orgID       string
		setupMock   func(*mocks.MockOrganizationRepository)
		assertErr   assert.ErrorAssertionFunc
		assertValue assert.ValueAssertionFunc
	}{
		{
			name:  "success - finds existing organization",
			orgID: orgID.String(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindOneByID(mock.Anything, orgID).
					Return(expectedOrg, nil).
					Once()
			},
			assertErr:   assert.NoError,
			assertValue: assert.NotNil,
		},
		{
			name:  "error - organization not found",
			orgID: nonExistingOrgID.String(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindOneByID(mock.Anything, nonExistingOrgID).
					Return(nil, errors.New("organization not found")).
					Once()
			},
			assertErr:   assert.Error,
			assertValue: assert.Nil,
		},
		{
			name:  "error - repository query failure",
			orgID: orgID.String(),
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindOneByID(mock.Anything, orgID).
					Return(nil, errors.New("database connection failed")).
					Once()
			},
			assertErr:   assert.Error,
			assertValue: assert.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockOrganizationRepository(t)
			tt.setupMock(repo)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			ctx := context.Background()

			uc, err := organization.NewUseCase(
				organization.UseCaseWithGenerator(uuid.New),
				organization.UseCaseWithLogger(logger),
				organization.UseCaseWithRepository(repo),
				organization.UseCaseWithUUIDParser(uuid.Parse),
			)
			require.NoError(t, err)

			result, err := uc.FindOneByID(ctx, tt.orgID)

			tt.assertErr(t, err)
			tt.assertValue(t, result)
		})
	}
}

func TestUseCaseFindAll(t *testing.T) {
	t.Parallel()

	expectedOrgs := []domain.Organization{
		{Identifier: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"), Name: "Org 1"},
		{Identifier: uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"), Name: "Org 2"},
	}

	tests := []struct {
		name        string
		limit       int
		offset      int
		setupMock   func(*mocks.MockOrganizationRepository)
		assertErr   assert.ErrorAssertionFunc
		assertValue assert.ValueAssertionFunc
	}{
		{
			name:   "success - returns multiple organizations",
			limit:  10,
			offset: 0,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindAll(mock.Anything, 10, 0).
					Return(expectedOrgs, nil).
					Once()
			},
			assertErr:   assert.NoError,
			assertValue: assert.NotNil,
		},
		{
			name:   "success - empty list when no organizations",
			limit:  10,
			offset: 0,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindAll(mock.Anything, 10, 0).
					Return([]domain.Organization{}, nil).
					Once()
			},
			assertErr:   assert.NoError,
			assertValue: assert.NotNil,
		},
		{
			name:   "success - pagination with offset",
			limit:  5,
			offset: 10,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindAll(mock.Anything, 5, 10).
					Return(expectedOrgs[:1], nil).
					Once()
			},
			assertErr:   assert.NoError,
			assertValue: assert.NotNil,
		},
		{
			name:   "error - repository query failure",
			limit:  10,
			offset: 0,
			setupMock: func(repo *mocks.MockOrganizationRepository) {
				repo.EXPECT().
					FindAll(mock.Anything, 10, 0).
					Return(nil, errors.New("database connection failed")).
					Once()
			},
			assertErr:   assert.Error,
			assertValue: assert.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockOrganizationRepository(t)
			tt.setupMock(repo)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			ctx := context.Background()

			uc, err := organization.NewUseCase(
				organization.UseCaseWithGenerator(uuid.New),
				organization.UseCaseWithLogger(logger),
				organization.UseCaseWithRepository(repo),
				organization.UseCaseWithUUIDParser(uuid.Parse),
			)
			require.NoError(t, err)

			result, err := uc.FindAll(ctx, tt.limit, tt.offset)

			tt.assertErr(t, err)
			tt.assertValue(t, result)
		})
	}
}

func TestUseCaseUpdate(t *testing.T) {
	t.Parallel()

	orgID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	nonExistingOrgID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
	type args struct {
		existingOrg *domain.Organization
		updateInput domain.UpdateOrganization
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(*mocks.MockOrganizationRepository, *domain.Organization)
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name: "success - updates organization name",
			args: args{
				existingOrg: &domain.Organization{
					Identifier: orgID,
					Name:       "Old Name",
				},
				updateInput: domain.UpdateOrganization{Name: "New Name"},
			},
			setupMock: func(repo *mocks.MockOrganizationRepository, existingOrg *domain.Organization) {
				repo.EXPECT().
					FindOneByID(mock.Anything, orgID).
					Return(existingOrg, nil).
					Once()
				repo.EXPECT().
					Update(mock.Anything, mock.MatchedBy(func(org *domain.Organization) bool {
						return org.Identifier.String() == orgID.String() && org.Name == "New Name"
					})).
					Return(nil).
					Once()
			},
			assertErr: assert.NoError,
		},
		{
			name: "error - repository update failure",
			args: args{
				existingOrg: &domain.Organization{
					Identifier: nonExistingOrgID,
				},
				updateInput: domain.UpdateOrganization{Name: "New Name"},
			},
			setupMock: func(repo *mocks.MockOrganizationRepository, org *domain.Organization) {
				repo.EXPECT().
					FindOneByID(mock.Anything, nonExistingOrgID).
					Return(nil, errors.New("organization not found")).
					Once()
			},
			assertErr: assert.Error,
		},
		{

			name: "error - repository update failure",
			args: args{
				existingOrg: &domain.Organization{
					Identifier: nonExistingOrgID,
				},
				updateInput: domain.UpdateOrganization{Name: "New Name"},
			},
			setupMock: func(repo *mocks.MockOrganizationRepository, org *domain.Organization) {
				repo.EXPECT().
					FindOneByID(mock.Anything, nonExistingOrgID).
					Return(org, nil).
					Once()
				repo.EXPECT().
					Update(mock.Anything, mock.Anything).
					Return(errors.New("database connection failed")).
					Once()
			},
			assertErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockOrganizationRepository(t)
			tt.setupMock(repo, tt.args.existingOrg)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			ctx := context.Background()

			uc, err := organization.NewUseCase(
				organization.UseCaseWithGenerator(uuid.New),
				organization.UseCaseWithLogger(logger),
				organization.UseCaseWithRepository(repo),
				organization.UseCaseWithUUIDParser(uuid.Parse),
			)
			require.NoError(t, err)

			err = uc.Update(ctx, tt.args.existingOrg.Identifier.String(), tt.args.updateInput)

			tt.assertErr(t, err)
		})
	}
}

func TestUseCaseDelete(t *testing.T) {
	t.Parallel()

	orgID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	nonExistingOrgID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

	type args struct {
		existingOrg *domain.Organization
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(*mocks.MockOrganizationRepository, *domain.Organization)
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name: "success - deletes organization",
			args: args{
				existingOrg: &domain.Organization{
					Identifier: orgID,
					Name:       "default",
				},
			},
			setupMock: func(repo *mocks.MockOrganizationRepository, org *domain.Organization) {
				repo.EXPECT().
					Delete(mock.Anything, org.Identifier).
					Return(nil).
					Once()
			},
			assertErr: assert.NoError,
		},
		{
			name: "error - organization not found",
			args: args{
				existingOrg: &domain.Organization{
					Identifier: nonExistingOrgID,
				},
			},
			setupMock: func(repo *mocks.MockOrganizationRepository, org *domain.Organization) {
				repo.EXPECT().
					Delete(mock.Anything, nonExistingOrgID).
					Return(errors.New("organization not found")).
					Once()
			},
			assertErr: assert.Error,
		},
		{
			name: "error - repository delete failure",
			args: args{
				existingOrg: &domain.Organization{
					Identifier: orgID,
				},
			},
			setupMock: func(repo *mocks.MockOrganizationRepository, org *domain.Organization) {
				repo.EXPECT().
					Delete(mock.Anything, org.Identifier).
					Return(errors.New("database constraint violation")).
					Once()
			},
			assertErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockOrganizationRepository(t)
			tt.setupMock(repo, tt.args.existingOrg)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			ctx := context.Background()

			uc, err := organization.NewUseCase(
				organization.UseCaseWithGenerator(uuid.New),
				organization.UseCaseWithLogger(logger),
				organization.UseCaseWithRepository(repo),
				organization.UseCaseWithUUIDParser(uuid.Parse),
			)
			require.NoError(t, err)

			err = uc.Delete(ctx, tt.args.existingOrg.Identifier.String())

			tt.assertErr(t, err)
		})
	}
}
