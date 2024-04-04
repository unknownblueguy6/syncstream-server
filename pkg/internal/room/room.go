package room

import (
	"errors"
	"fmt"
	"math/rand"
	"syncstream-server/pkg/internal/stream"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var Manager = &RoomManager{
	Map:       make(map[RoomCode]*Room),
	Events:    make(chan *Event),
	Users:     make(map[uuid.UUID]*RoomUser),
	UserIDMap: make(map[uuid.UUID]uuid.UUID),
}

const EVENT_BUFFER_SIZE = 10

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

func (manager *RoomManager) AddRoom(id uuid.UUID, url string, streamState stream.StreamState, streamElement stream.StreamElement) (RoomCode, error) {
	code := generateRoomCode()
	_, ok := manager.Map[code]
	for ok {
		code = generateRoomCode()
		_, ok = manager.Map[code]
	}
	manager.Map[code] = NewRoom(id, code, url, streamState, streamElement)
	return code, nil
}

func (manager *RoomManager) Run() {
	fmt.Println("Starting Room Manager")
	for {
		select {
		case event, ok := <-manager.Events:
			if !ok {
				fmt.Println("RoomManger Events Channel closed")
				return
			}
			fmt.Println(*event)
			var code RoomCode
			var room *Room

			switch event.Type {
			case USER_JOIN:
				user := event.Data["user"].(*RoomUser)

				manager.Users[event.SourceID] = user
				manager.UserIDMap[event.SourceID] = uuid.New()

				code = user.Code
				room = manager.Map[code]
				room.Users[event.SourceID] = true

				event.Data = nil

				user.Events <- manager.Map[code].ToEvent(manager.UserIDMap[event.SourceID])

			case USER_LEFT:
				code = manager.Users[event.SourceID].Code
				room = manager.Map[code]

				room.Users[event.SourceID] = false
				delete(manager.UserIDMap, event.SourceID)
				delete(manager.Users, event.SourceID)

			default:
				code = manager.Users[event.SourceID].Code
				room = manager.Map[code]
			}

			sourceID := event.SourceID
			event.SourceID = manager.UserIDMap[event.SourceID]

			fmt.Println(sourceID, event.SourceID)

			for userID := range room.Users {
				if userID != sourceID {
					manager.Users[userID].Events <- event
				}
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
