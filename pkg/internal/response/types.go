package response

import (
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
)

type InitResponseBody struct {
	ID uuid.UUID `json:"id"`
	// bearerAuth string
}

type CreateResponseBody struct {
	Code room.RoomCode `json:"code"`
}

type JoinResponseBody struct {
	URL           string               `json:"url"`
	StreamState   stream.StreamState   `json:"streamState"`
	StreamElement stream.StreamElement `json:"streamElement"`
	Timestamp     time.Time            `json:"timestamp"`
}

type DeleteResponseBody struct {
	Success bool `json:"success"`
}
