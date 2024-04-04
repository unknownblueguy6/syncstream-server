package request

import (
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
)

type InitRequestBody struct {
	// ID room.UUID `json:"id"`
}

type CreateRequestBody struct {
	ID            uuid.UUID            `json:"id"`
	URL           string               `json:"url"`
	StreamState   stream.StreamState   `json:"streamState"`
	StreamElement stream.StreamElement `json:"streamElement"`
	Timestamp     time.Time            `json:"timestamp"`
}

type JoinTokenRequestBody struct {
	ID   uuid.UUID     `json:"id"`
	Code room.RoomCode `json:"code"`
}

type DeleteRequestBody struct {
	ID   uuid.UUID     `json:"id"`
	Code room.RoomCode `json:"code"`
}
