package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hoppermq/streamly/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestUsecaseCreate(t *testing.T) {
	t.Parallel()

	type args struct {
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(*mocks.MockUserRepository)
		asserErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success - create user",
			args: args{},
			setupMock: func(m *mocks.MockUserRepository) {
			},
			asserErr: assert.NoError,
		},
	}
}
