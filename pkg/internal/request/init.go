package request

import (
	"encoding/json"
	"fmt"
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
	if r.Method == http.MethodPost {
		var reqBody InitRequestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resBody := InitResponseBody{ID: uuid.New()}
		fmt.Println(*r, resBody)
		err = json.NewEncoder(w).Encode(resBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
