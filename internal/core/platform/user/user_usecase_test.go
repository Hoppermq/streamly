package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsecaseCreate(t *testing.T) {
	t.Parallel()

	type args struct {
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(any2 any)
		asserErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success - create user",
			args: args{},
			setupMock: func(m any) {
			},
			asserErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})
	}
}
