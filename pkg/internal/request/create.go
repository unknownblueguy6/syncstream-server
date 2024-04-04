package request

import (
	"encoding/json"
	"fmt"
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
	if r.Method == http.MethodPost {
		var reqBody CreateRequestBody

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		code, err := room.Manager.AddRoom(reqBody.ID, reqBody.URL, reqBody.StreamState, reqBody.StreamElement)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		resBody := CreateResponseBody{Code: code}
		fmt.Println(*r, resBody)
		err = json.NewEncoder(w).Encode(resBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
