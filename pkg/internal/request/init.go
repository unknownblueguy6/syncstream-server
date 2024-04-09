package request

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type InitRequestBody struct {
	// ID room.UUID `json:"id"`
}

type InitResponseBody struct {
	ID uuid.UUID `json:"id"`
	// bearerAuth string
}

func InitHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody InitRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Error("POST /init", "error", err.Error(), "body", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("POST /init", "req_body", reqBody)

	resBody := InitResponseBody{ID: uuid.New()}
	slog.Debug("POST /init", "res_body", resBody)

	err = json.NewEncoder(w).Encode(resBody)
	if err != nil {
		slog.Error("POST /init", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
