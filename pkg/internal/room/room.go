package room

import (
	"errors"
	"math/rand"
	"syncstream-server/pkg/internal/stream"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type RoomCode string

const ROOMCODE_LENGTH = 6
const ROOMCODE_CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Room struct {
	Code          RoomCode
	creatorID     uuid.UUID
	Users         map[uuid.UUID]bool
	Empty         bool
	URL           string
	StreamState   stream.StreamState
	StreamElement stream.StreamElement
}

func NewRoom(id uuid.UUID, code RoomCode, url string, ss stream.StreamState, se stream.StreamElement) *Room {
	return &Room{
		Code:          code,
		creatorID:     id,
		Users:         make(map[uuid.UUID]bool),
		Empty:         true,
		URL:           url,
		StreamState:   ss,
		StreamElement: se,
	}
}

func (room *Room) ToEvent(id uuid.UUID) *Event {
	users := []string{}
	for id := range room.Users {
		users = append(users, Manager.UserIDMap[id].String())
	}
	return &Event{
		SourceID:  id,
		Timestamp: time.Now(),
		Type:      ROOM_STATE,
		Data: map[string]any{
			"url":           room.URL,
			"streamState":   room.StreamState,
			"streamElement": room.StreamElement,
			"users":         users,
		},
	}
}

func generateRoomCode() RoomCode {
	code := make([]byte, ROOMCODE_LENGTH)
	for i := range code {
		code[i] = ROOMCODE_CHARSET[rand.Intn(len(ROOMCODE_CHARSET))]
	}

	return RoomCode(code)
}

func ParseRoomCode(code string) (RoomCode, error) {
	if len(code) != ROOMCODE_LENGTH {
		return RoomCode(""), errors.New("invalid RoomCode length")
	}
	for _, char := range code {
		if !unicode.IsUpper(char) {
			return RoomCode(""), errors.New("invalid RoomCode character")
		}
	}
	return RoomCode(code), nil
}
