package request

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"syncstream-server/pkg/internal/valid8r"

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

	if errs := valid8r.Validator.Struct(&reqBody); errs != nil {
		for _, err := range errs.(valid8r.ValidationErrors) {
			slog.Error(err.Error())
		}
		http.Error(w, errs.Error(), http.StatusBadRequest)
		return
	}

	resBody := InitResponseBody{ID: uuid.New()}
	slog.Debug("POST /init", "res_body", resBody)

	err = json.NewEncoder(w).Encode(resBody)
	if err != nil {
		slog.Error("POST /init", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
