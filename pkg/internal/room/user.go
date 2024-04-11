package room

import (
	"log/slog"
	"syncstream-server/pkg/internal/valid8r"
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
				slog.Debug(user.id.String()+" ReceiveEventsFromClient() USER_LEFT", "code", user.Code)
				user.Manager.Events <- &Event{SourceID: user.id, Timestamp: time.Now(), Type: USER_LEFT, Data: nil}
			} else {
				slog.Error(err.Error())
			}

			return
		}

		slog.Debug(user.id.String()+" ReceiveEventsFromClient()", "receivedEventFromClient", *event)

		if errs := valid8r.Validator.Struct(event); errs != nil {
			for _, err := range errs.(valid8r.ValidationErrors) {
				slog.Error(err.Error())
			}
			return
		}

		if event.Type != ZERO && event.IsValid(user.id) {
			slog.Debug(user.id.String()+" ReceiveEventsFromClient()", "validEvent", *event)
			user.Manager.Events <- event
		} else {
			slog.Debug(user.id.String()+" ReceiveEventsFromClient()", "invalidEvent", *event)
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
				slog.Error(user.id.String() + " Events Channel closed")
				return
			}
			err := user.Conn.WriteJSON(event)
			if err != nil {
				slog.Error(user.id.String()+" SendEventsToClient()", "error", err.Error(), "event", *event)
				continue
			}
			slog.Debug(user.id.String()+" SendEventsToClient()", "event", *event)
		default:
			time.Sleep(100 * time.Millisecond)
		}

	}
}
