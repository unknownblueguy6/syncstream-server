package room

import (
	"log/slog"
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
)

type EphemeralTokenData struct {
	ID         uuid.UUID
	Code       RoomCode
	ExpiryTime time.Time
}

type RoomManager struct {
	Map       map[RoomCode]*Room
	Events    chan *Event
	Users     map[uuid.UUID]*RoomUser
	UserIDMap map[uuid.UUID]uuid.UUID
	Tokens    map[uuid.UUID]EphemeralTokenData
}

// TODO : update Room state every second

var Manager = &RoomManager{
	Map:       make(map[RoomCode]*Room),
	Events:    make(chan *Event),
	Users:     make(map[uuid.UUID]*RoomUser),
	UserIDMap: make(map[uuid.UUID]uuid.UUID),
	Tokens:    make(map[uuid.UUID]EphemeralTokenData),
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
	slog.Info("Starting Room Manager")
	for {
		select {
		case event, ok := <-manager.Events:
			if !ok {
				slog.Error("RoomManger Events Channel closed")
				return
			}
			slog.Debug("manager.Run()", "receivedEvent", *event)
			var code RoomCode
			var room *Room

			switch event.Type {
			case USER_JOIN:
				user := event.Data["user"].(*RoomUser)

				manager.Users[event.SourceID] = user
				if _, ok := manager.UserIDMap[event.SourceID]; !ok {
					manager.UserIDMap[event.SourceID] = uuid.New()
				}

				code = user.Code
				room = manager.Map[code]
				room.Users[event.SourceID] = true

				event.Data = nil

				slog.Debug("manager.Run() USER_JOIN", "RoomStateEvent", *manager.Map[code].ToEvent(manager.UserIDMap[event.SourceID]))
				user.Events <- manager.Map[code].ToEvent(manager.UserIDMap[event.SourceID])

			case USER_LEFT:
				code = manager.Users[event.SourceID].Code
				room = manager.Map[code]

				room.Users[event.SourceID] = false
				delete(manager.Users, event.SourceID)

			default:
				code = manager.Users[event.SourceID].Code
				room = manager.Map[code]
			}

			sourceID := event.SourceID
			event.SourceID = manager.UserIDMap[sourceID]
			slog.Debug("manager.Run()", "orig_id", sourceID, "mapped_id", event.SourceID)
			slog.Debug("manager.Run()", "sentEvent", *event)
			for userID := range room.Users {
				if userID != sourceID {
					slog.Debug("manager.Run() "+sourceID.String()+" Event", "destinationID", userID)
					manager.Users[userID].Events <- event
				}
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
