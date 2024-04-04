package room

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomUser struct {
	id      uuid.UUID
	Code    RoomCode
	Events  chan *Event
	Conn    *websocket.Conn
	Manager *RoomManager
}

func NewRoomUser(id uuid.UUID, code RoomCode, conn *websocket.Conn, manager *RoomManager) *RoomUser {
	return &RoomUser{
		id:      id,
		Code:    code,
		Conn:    conn,
		Manager: manager,
		Events:  make(chan *Event, EVENT_BUFFER_SIZE),
	}
}

func (user *RoomUser) ReceiveEventsFromClient() {
	defer func() {
		user.Conn.Close()
	}()
	for {
		event := &Event{SourceID: uuid.UUID{}, Timestamp: time.Time{}, Type: ZERO, Data: nil}
		err := user.Conn.ReadJSON(event)
		if err != nil {

			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				fmt.Println(err)
				user.Manager.Events <- &Event{SourceID: user.id, Timestamp: time.Now(), Type: USER_LEFT, Data: nil}
			} else {
				fmt.Println("could not parse to Event Type", err)
			}

			return
		}
		fmt.Println(*event)
		if event.Type != ZERO && event.IsValid(user.id) {
			user.Manager.Events <- event
		}
	}
}

func (user *RoomUser) SendEventsToClient() {
	// ticker := time.NewTicker()
	defer func() {
		user.Conn.Close()
	}()
	for {
		select {
		case event, ok := <-user.Events:
			if !ok {
				continue
			}
			err := user.Conn.WriteJSON(event)
			if err != nil {
				continue
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}

	}
}
