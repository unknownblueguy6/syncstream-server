package room

import (
	"fmt"
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
)

type RoomManager struct {
	Map       map[RoomCode]*Room
	Events    chan *Event
	Users     map[uuid.UUID]*RoomUser
	UserIDMap map[uuid.UUID]uuid.UUID
}

// TODO : make mapped ID persist across disconnects, i.e., a given uuid is always mapped to the same uuid.
// TODO : update Room state every second

var Manager = &RoomManager{
	Map:       make(map[RoomCode]*Room),
	Events:    make(chan *Event),
	Users:     make(map[uuid.UUID]*RoomUser),
	UserIDMap: make(map[uuid.UUID]uuid.UUID),
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
