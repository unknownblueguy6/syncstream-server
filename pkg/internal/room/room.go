package room

import (
	"errors"
	"log/slog"
	"math/rand"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type RoomCode string

type StreamState struct {
	CurrentTime  float64 `json:"currentTime"  validate:"numeric"`
	Paused       bool    `json:"paused" validate:"boolean"`
	PlaybackRate float32 `json:"playbackRate" validate:"required,numeric"`
}

type StreamElement any

const ROOMCODE_LENGTH = 6
const ROOMCODE_CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Room struct {
	Code          RoomCode
	creatorID     uuid.UUID
	Users         map[uuid.UUID]bool
	Empty         bool
	URL           string
	StreamState   StreamState
	StreamElement StreamElement
	lastUpdate    time.Time
}

func NewRoom(id uuid.UUID, code RoomCode, url string, ss StreamState, se StreamElement) *Room {
	return &Room{
		Code:          code,
		creatorID:     id,
		Users:         make(map[uuid.UUID]bool),
		Empty:         true,
		URL:           url,
		StreamState:   ss,
		StreamElement: se,
		lastUpdate:    time.Now().UTC(),
	}
}

func (room *Room) ToEvent(id uuid.UUID) *Event {
	users := []string{}
	for id := range room.Users {
		users = append(users, Manager.UserIDMap[id].String())
	}
	return &Event{
		SourceID:  id,
		Timestamp: room.lastUpdate,
		Type:      ROOM_STATE,
		Data: map[string]any{
			"url":           room.URL,
			"streamState":   room.StreamState,
			"streamElement": room.StreamElement,
			"users":         users,
		},
	}
}

func (room *Room) UpdateStream() {
	if !room.StreamState.Paused {
		t := time.Now().UTC()
		timeDelta := t.Sub(room.lastUpdate).Seconds()
		room.StreamState.CurrentTime += timeDelta * float64(room.StreamState.PlaybackRate)
		room.lastUpdate = t
	}
	slog.Debug(string(room.Code)+" UpdateStream()", "streamState", room.StreamState)
}

func (room *Room) UpdateStreamEvent(event *Event) {
	slog.Debug(string(room.Code)+" UpdateStreamEvent()", "event", event)
	room.StreamState = event.GetStreamState()
	room.lastUpdate = event.Timestamp
	room.UpdateStream()
}

func generateRoomCode() RoomCode {
	code := make([]byte, ROOMCODE_LENGTH)
	for i := range code {
		code[i] = ROOMCODE_CHARSET[rand.Intn(len(ROOMCODE_CHARSET))]
	}

	return RoomCode(code)
}

func ParseRoomCode(code RoomCode) (RoomCode, error) {
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
