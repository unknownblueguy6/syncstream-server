package request

import (
	"fmt"
	"net/http"
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO : make this route return ephemeral token for websocket connection
// TODO : use Get() method on url.Query() to get query params

type JoinTokenRequestBody struct {
	ID   uuid.UUID     `json:"id"`
	Code room.RoomCode `json:"code"`
}

type JoinResponseBody struct {
	URL           string               `json:"url"`
	StreamState   stream.StreamState   `json:"streamState"`
	StreamElement stream.StreamElement `json:"streamElement"`
	Timestamp     time.Time            `json:"timestamp"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JoinHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(*r)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	queries := r.URL.Query()

	_, codeOK := queries["code"]
	_, idOK := queries["id"]

	if !codeOK || !idOK {
		conn.Close()
		return
	}

	code, codeErr := room.ParseRoomCode(queries["code"][0])
	id, idErr := uuid.Parse(queries["id"][0])

	if codeErr != nil || idErr != nil {
		conn.Close()
		return
	}

	_, roomOK := room.Manager.Map[code]
	if !roomOK {
		conn.Close()
		return
	}

	user := room.NewRoomUser(id, code, conn, room.Manager)
	room.Manager.Events <- &room.Event{SourceID: id, Timestamp: time.Now(), Type: room.USER_JOIN, Data: map[string](any){"user": user}}

	go user.ReceiveEventsFromClient()
	go user.SendEventsToClient()

}
