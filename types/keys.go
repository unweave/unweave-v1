package types

import (
	"github.com/google/uuid"
)

// Context Keys

const (
	ContextKeyUser = "user"
)

type UserContext struct {
	ID uuid.UUID `json:"id"`
}
