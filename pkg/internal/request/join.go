package request

import (
	"fmt"
	"log/slog"
	"net/http"
	"syncstream-server/pkg/internal/room"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func JoinHandler(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	token, err := uuid.Parse(queries.Get("token"))
	if err != nil {
		slog.Debug("ws /join?{token}", "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenData, tokenOK := room.Manager.Tokens[token]
	if !tokenOK {
		slog.Debug("ws /join?{token}", "error", "Invalid Token")
		http.Error(w, "Invalid Token", http.StatusForbidden)
		return
	}

	if tokenData.ExpiryTime.Before(time.Now().UTC()) {
		slog.Debug("ws /join?{token}", "error", "Expired Token")
		http.Error(w, "Expired Token", http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	user := room.NewRoomUser(tokenData.ID, tokenData.Code, conn, room.Manager)
	room.Manager.Events <- &room.Event{SourceID: tokenData.ID, Timestamp: time.Now(), Type: room.USER_JOIN, Data: map[string](any){"user": user}}

	go user.ReceiveEventsFromClient()
	go user.SendEventsToClient()
}
