package room

import (
	"time"

	"github.com/google/uuid"
)

type EventType int

const (
	ZERO EventType = iota
	PLAY
	PAUSE
	SEEK
	USER_JOIN
	USER_LEFT
	MESSAGE
	ROOM_STATE
)

const EVENT_BUFFER_SIZE = 10

type Event struct {
	SourceID  uuid.UUID      `json:"sourceID" validate:"required,uuid"`
	Timestamp time.Time      `json:"timestamp" validate:"required"`
	Type      EventType      `json:"type" validate:"required,number"`
	Data      map[string]any `json:"data"`
}

// TODO add validation for the Data field for all possible events

func (e *Event) IsValid(id uuid.UUID) bool {
	switch {
	case !(e.Type > ZERO && e.Type <= MESSAGE):
		return false
	case id != e.SourceID:
		return false
	// case e.Data != nil:
	// 	keys := []string{}
	// 	for k, _ := range e.Data {
	// 		keys = append(keys, k)
	// 	}
	// 	slices.SortFunc[]()

	default:
		return true
	}
}
