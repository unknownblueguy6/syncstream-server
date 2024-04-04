package request

import (
	"syncstream-server/pkg/internal/room"

	"github.com/google/uuid"
)

type DeleteRequestBody struct {
	ID   uuid.UUID     `json:"id"`
	Code room.RoomCode `json:"code"`
}

type DeleteResponseBody struct {
	Success bool `json:"success"`
}
