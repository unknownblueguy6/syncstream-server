package request

import (
	"encoding/json"
	"log/slog"

	"net/http"
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/stream"
	"time"

	"github.com/google/uuid"
)

type CreateRequestBody struct {
	ID            uuid.UUID            `json:"id"`
	URL           string               `json:"url"`
	StreamState   stream.StreamState   `json:"streamState"`
	StreamElement stream.StreamElement `json:"streamElement"`
	Timestamp     time.Time            `json:"timestamp"`
}

type CreateResponseBody struct {
	Code room.RoomCode `json:"code"`
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody CreateRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Error("POST /create", "error", err.Error(), "body", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("POST /create", "req_body", reqBody)

	code, err := room.Manager.AddRoom(reqBody.ID, reqBody.URL, reqBody.StreamState, reqBody.StreamElement)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resBody := CreateResponseBody{Code: code}
	slog.Debug("POST /create", "res_body", resBody)

	err = json.NewEncoder(w).Encode(resBody)
	if err != nil {
		slog.Error("POST /create", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
