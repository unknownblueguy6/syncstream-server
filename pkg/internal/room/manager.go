package room

import (
	"log/slog"
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

var Manager = &RoomManager{
	Map:       make(map[RoomCode]*Room),
	Events:    make(chan *Event),
	Users:     make(map[uuid.UUID]*RoomUser),
	UserIDMap: make(map[uuid.UUID]uuid.UUID),
	Tokens:    make(map[uuid.UUID]EphemeralTokenData),
}

func (manager *RoomManager) AddRoom(id uuid.UUID, url string, streamState StreamState, streamElement StreamElement) (RoomCode, error) {
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
			var sourceID = event.SourceID
			if _, ok := manager.UserIDMap[sourceID]; !ok {
				manager.UserIDMap[sourceID] = uuid.New()
			}
			var mappedID = manager.UserIDMap[sourceID]

			// get the room code and room from the event
			switch event.Type {
			case USER_JOIN:
				user := event.Data["user"].(*RoomUser)

				manager.Users[sourceID] = user

				code = user.Code
				room = manager.Map[code]
				room.Users[sourceID] = true

				event.Data = nil

				room.UpdateStream()

				roomStateEvent := manager.Map[code].ToEvent(mappedID)
				slog.Debug("manager.Run() USER_JOIN", "RoomStateEvent", *roomStateEvent)
				user.Events <- roomStateEvent

			case USER_LEFT:
				code = manager.Users[sourceID].Code
				room = manager.Map[code]

				delete(room.Users, sourceID)
				delete(manager.Users, sourceID)

			default:
				code = manager.Users[sourceID].Code
				room = manager.Map[code]
			}

			event.SourceID = mappedID

			slog.Debug("manager.Run()", "orig_id", sourceID, "mapped_id", mappedID)
			slog.Debug("manager.Run()", "sentEvent", *event)

			if event.IsStreamEvent() {
				room.UpdateStreamEvent(event)
			}

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
