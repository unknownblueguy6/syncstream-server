package room

import (
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

type Event struct {
	SourceID  uuid.UUID      `json:"sourceID"`
	Timestamp time.Time      `json:"timestamp"`
	Type      EventType      `json:"type"`
	Data      map[string]any `json:"data"`
}

type RoomCode string

const ROOMCODE_LENGTH = 6
const ROOMCODE_CHARSET = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type RoomUser struct {
	id      uuid.UUID
	Code    RoomCode
	Events  chan *Event
	Conn    *websocket.Conn
	Manager *RoomManager
}

type RoomManager struct {
	Map       map[RoomCode]*Room
	Events    chan *Event
	Users     map[uuid.UUID]*RoomUser
	UserIDMap map[uuid.UUID]uuid.UUID
}

type Room struct {
	Code          RoomCode
	creatorID     uuid.UUID
	Users         map[uuid.UUID]bool
	Empty         bool
	URL           string
	StreamState   stream.StreamState
	StreamElement stream.StreamElement
}
