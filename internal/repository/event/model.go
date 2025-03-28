package event

import (
	"time"

	"github.com/google/uuid"
)

//go:generate eos generator repository --type=Event --common_path=pkg/common
type Event struct {
	ID        uuid.UUID
	Payload   []byte
	CreatedAt time.Time
}
