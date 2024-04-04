package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syncstream-server/pkg/internal/response"

	"github.com/google/uuid"
)

func InitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqBody InitRequestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resBody := response.InitResponseBody{ID: uuid.New()}
		fmt.Println(*r, resBody)
		err = json.NewEncoder(w).Encode(resBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
