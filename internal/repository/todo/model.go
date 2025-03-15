package todo

import (
	"time"

	"github.com/google/uuid"
)

//go:generate eos generator repository --type=Task --common_path=pkg/common
type Task struct {
	ID uuid.UUID

	Name        string
	Description *string

	CreatedBy string
	Doer      string
	Done      bool

	Repeatable  bool
	RepeatAfter *int

	DoBefore  *time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
}
