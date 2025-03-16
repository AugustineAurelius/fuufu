package user

import (
	"time"

	"github.com/google/uuid"
)

//go:generate eos generator repository --type=User --common_path=pkg/common
type User struct {
	ID uuid.UUID

	Name           string
	Email          string
	HashedPassword string

	CreatedAt time.Time
	UpdatedAt *time.Time
}
