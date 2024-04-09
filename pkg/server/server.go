package server

import (
	"flag"
	"fmt"
	"net/http"

	"syncstream-server/pkg/internal/request"
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/stream"

	"github.com/google/uuid"
)

var addr = flag.String("addr", "localhost:8080", "Address of Server")

func Run() {
	// TODO add error logging
	// TODO make ci/cd pipeline work
	// TODO write tests

	id, err := uuid.NewUUID()
	if err != nil {
		return
	}
	fmt.Println(id)
	room.Manager.Map["ABCDEF"] = room.NewRoom(id, "ABCDEF", "http://example.com", stream.StreamState{CurrentTime: 0.0, Paused: false, PlaybackRate: 1.0}, nil)
	go room.Manager.Run()

	fmt.Println("Starting Server")
	http.HandleFunc("/init", request.InitHandler)
	http.HandleFunc("/create", request.CreateHandler)
	http.HandleFunc("/join", request.JoinHandler)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		return
	}
}
