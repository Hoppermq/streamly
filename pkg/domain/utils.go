package domain

import "github.com/google/uuid"

type Generator func() uuid.UUID
type UUIDParser func(string) (uuid.UUID, error)
