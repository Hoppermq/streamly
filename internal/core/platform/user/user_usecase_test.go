package user_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/internal/core/platform/user"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/hoppermq/streamly/pkg/domain/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUseCaseCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		userInput     *domain.CreateUser
		generatedUUID uuid.UUID
	}

	tests := []struct {
		name          string
		args          args
		setupMock     func(repository *mocks.MockUserRepository, input *domain.CreateUser)
		asserErr      assert.ErrorAssertionFunc
		expectedError error
	}{
		{
			name: "success - create user",
			args: args{
				userInput: &domain.CreateUser{
					UserName:     "test",
					FirstName:    "test",
					LastName:     "test",
					PrimaryEmail: "test@tester.com",
					ZitadelID:    "354117604350126857",
					Role:         "user",
				},
				generatedUUID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
			},
			setupMock: func(r *mocks.MockUserRepository, input *domain.CreateUser) {
				r.EXPECT().Create(mock.Anything, mock.MatchedBy(func(user *domain.User) bool {
					return user.ZitadelID == input.ZitadelID
				})).Return(nil).Once()
			},
			asserErr: assert.NoError,
		}, {
			name: "error - repository fails to insert",
			args: args{
				userInput: &domain.CreateUser{
					UserName:     "test",
					FirstName:    "test",
					LastName:     "test",
					PrimaryEmail: "test@tester",
					ZitadelID:    "354117604350126857",
					Role:         "user",
				},
			},
			setupMock: func(r *mocks.MockUserRepository, input *domain.CreateUser) {
				r.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("error"))
			},
			asserErr: assert.Error,
		}, {
			name: "failure - nil user input",
			args: args{
				userInput:     nil,
				generatedUUID: uuid.New(),
			},
			setupMock: func(r *mocks.MockUserRepository, input *domain.CreateUser) {},
			asserErr:  assert.Error,
		}, {
			name: "failure - duplicate user input",
			args: args{
				userInput: &domain.CreateUser{
					UserName:     "test",
					FirstName:    "test",
					LastName:     "test",
					PrimaryEmail: "test@tester",
					ZitadelID:    "354117604350126857",
					Role:         "user",
				},
				generatedUUID: uuid.New(),
			},
			setupMock: func(r *mocks.MockUserRepository, input *domain.CreateUser) {
				r.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(errors.New("pq: duplicate key value violates unique constraint")).Once()
			},
			asserErr: assert.Error,
		},
		{
			name: "failure - context cancellation",
			args: args{
				userInput: &domain.CreateUser{
					UserName:     "test",
					FirstName:    "test",
					LastName:     "test",
					PrimaryEmail: "test@tester",
					ZitadelID:    "354117604350126857",
					Role:         "user",
				},
			},
			setupMock: func(r *mocks.MockUserRepository, input *domain.CreateUser) {
				r.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(context.Canceled).Once()
			},
			asserErr: assert.Error,
		}, {
			name: "failure - missing zitadel ID",
			args: args{
				userInput: &domain.CreateUser{
					UserName:     "test",
					FirstName:    "test",
					LastName:     "test",
					PrimaryEmail: "test@tester",
					ZitadelID:    "",
					Role:         "user",
				},
			},
			setupMock: func(r *mocks.MockUserRepository, input *domain.CreateUser) {
				r.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).
					Return(errors.New("missing Zitadel ID")).Once()
			},
			asserErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewMockUserRepository(t)
			authRepo := mocks.NewMockAuthRepository(t)
			tt.setupMock(repo, tt.args.userInput)

			client := mocks.NewMockClient(t)

			logger := slog.New(slog.DiscardHandler)
			ctx := context.Background()

			uc, err := user.NewUseCase(
				user.UseCaseWithLogger(logger),
				user.UseCaseWithGenerator(func() uuid.UUID { return tt.args.generatedUUID }),
				user.UseCaseWithUUIDParser(uuid.Parse),
				user.UseCaseWithUserRepository(repo),
				user.UseCaseWithAuthRepository(authRepo),
				user.UseCaseWithZitadelAPI(client),
			)

			require.NoError(t, err)
			err = uc.Create(ctx, tt.args.userInput)
			tt.asserErr(t, err)
		})
	}
}

func TestUseCase_CreateFromEvent(t *testing.T) {
	t.Parallel()

	generatedUUID := uuid.MustParse("123e4567-e89b-12d3-a456-42661417400e")
	type args struct {
		event         *domain.ZitadelEventUserCreated
		generatedUUID uuid.UUID
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(*mocks.MockUserRepository, *mocks.MockClient)
		asserErr  assert.ErrorAssertionFunc
	}{{
		name: "success - create user from event",
		args: args{
			event: &domain.ZitadelEventUserCreated{
				InstanceID:     "123435464",
				OrganizationID: "123435464",
				UserID:         "123435464",
				Request: domain.ZitadelEventUserCreatedRequest{
					Email:        domain.ZitadelEmail{},
					Profile:      domain.ZitadelProfile{},
					Organization: domain.ZitadelOrganization{},
				},
			},
			generatedUUID: generatedUUID,
		},
		setupMock: func(ur *mocks.MockUserRepository, zr *mocks.MockClient) {
			zr.EXPECT().
				GetUserByUserName(mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("string")).
				Return(&domain.User{
					UserName:     "test",
					PrimaryEmail: "test",
					FirstName:    "test",
					LastName:     "test",
					ZitadelID:    "123435464",
				}, nil).
				Once()
			ur.EXPECT().Create(mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil).Once()
		},
		asserErr: assert.NoError,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockUserRepository(t)
			authRepo := mocks.NewMockAuthRepository(t)
			client := mocks.NewMockClient(t)

			tt.setupMock(repo, client)

			logger := slog.New(slog.DiscardHandler)
			ctx := context.Background()

			uc, err := user.NewUseCase(
				user.UseCaseWithLogger(logger),
				user.UseCaseWithGenerator(func() uuid.UUID { return tt.args.generatedUUID }),
				user.UseCaseWithUUIDParser(uuid.Parse),
				user.UseCaseWithUserRepository(repo),
				user.UseCaseWithAuthRepository(authRepo),
				user.UseCaseWithZitadelAPI(client),
			)
			require.NoError(t, err)
			err = uc.CreateFromEvent(ctx, tt.args.event)
			tt.asserErr(t, err)
		})
	}
}
