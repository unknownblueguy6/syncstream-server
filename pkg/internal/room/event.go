package room

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type EventType int

const (
	ZERO EventType = iota
	USER_JOIN
	USER_LEFT
	ROOM_STATE
	PLAY
	PAUSE
	SEEK
	MESSAGE
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
	case !(e.Type >= PLAY && e.Type <= MESSAGE):
		return false
	case id != e.SourceID:
		return false
	case e.Type >= PLAY && e.Type <= SEEK:
		_, ok := e.Data["streamState"].(map[string]any)
		return ok
	}
	return true
}

func (e *Event) IsStreamEvent() bool {
	return e.Type >= PLAY && e.Type <= SEEK
}

func (e *Event) GetStreamState() StreamState {
	slog.Debug("Converting event to stream state", "event", e)
	ss, _ := e.Data["streamState"].(map[string]any)
	slog.Debug("ss", "ss", ss)
	return StreamState{
		CurrentTime:  ss["currentTime"].(float64),
		Paused:       ss["paused"].(bool),
		PlaybackRate: float32(ss["playbackRate"].(float64)),
	}
}
