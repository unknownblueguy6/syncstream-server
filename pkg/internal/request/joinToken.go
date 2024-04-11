package request

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/valid8r"
	"time"

	"github.com/google/uuid"
)

type JoinTokenRequestBody struct {
	ID   uuid.UUID     `json:"id" validate:"required,uuid"`
	Code room.RoomCode `json:"code" validate:"required,alpha,uppercase,len=6"`
}

type JoinTokenResponseBody struct {
	Token      uuid.UUID `json:"token"`
	ExpiryTime time.Time `json:"expiryTime"`
}

func JoinTokenHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody JoinTokenRequestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Error("POST /joinToken", "error", err.Error(), "body", r.Body)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("POST /joinToken", "req_body", reqBody)

	if errs := valid8r.Validator.Struct(&reqBody); errs != nil {
		for _, err := range errs.(valid8r.ValidationErrors) {
			slog.Error(err.Error())
		}
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	code := reqBody.Code
	_, roomOK := room.Manager.Map[code]
	if !roomOK {
		slog.Error("POST /joinToken", "error", string(code)+" has not been created.")
		http.Error(w, string(code)+" has not been created.", http.StatusNotFound)
		return
	}

	token := uuid.New()
	expiryTime := time.Now().UTC().Add(6 * time.Hour)
	room.Manager.Tokens[token] = room.EphemeralTokenData{
		ID:         reqBody.ID,
		Code:       reqBody.Code,
		ExpiryTime: expiryTime,
	}

	resBody := JoinTokenResponseBody{
		Token:      token,
		ExpiryTime: expiryTime,
	}
	slog.Debug("POST /joinToken", "res_body", resBody)

	err = json.NewEncoder(w).Encode(resBody)
	if err != nil {
		slog.Error("POST /joinToken", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
