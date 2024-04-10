package request

import (
	"bytes"
	"encoding/json"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

// func check
// func (resBody *InitResponseBody) UnmarshalJSON(data []byte) error {
// 	var aux struct {
// 		ID string `json:"id"`
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	parsedID, err := uuid.Parse(aux.ID)
// 	if err != nil {
// 		return err
// 	}
// 	resBody.ID = parsedID
// 	return nil
// }

func initReqBodyToRequest(reqBody InitRequestBody) *http.Request {
	reqBodyBytes, _ := json.Marshal(reqBody)
	return httptest.NewRequest(http.MethodPost, "/init", bytes.NewBuffer(reqBodyBytes))
}

func TestInitHandler(t *testing.T) {
	tests := []struct {
		name  string
		req   *http.Request
		_func func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Correct Request",
			req:  initReqBodyToRequest(InitRequestBody{}),
			_func: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
				}
				resBody := InitResponseBody{ID: uuid.UUID{}}
				err := json.NewDecoder(rec.Body).Decode(&resBody)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}
			},
		},
		{
			name: "Empty Request",
			req:  httptest.NewRequest(http.MethodPost, "/init", nil),
			_func: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusBadRequest {
					t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rec.Code)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			InitHandler(rec, test.req)
			test._func(t, rec)
		})
	}
}
