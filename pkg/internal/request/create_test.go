package request

import (
	"bytes"
	"encoding/json"
	"syncstream-server/pkg/internal/room"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

func createReqBodyToRequest(reqBody CreateRequestBody) *http.Request {
	reqBodyBytes, _ := json.Marshal(reqBody)
	return httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(reqBodyBytes))
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name  string
		req   *http.Request
		_func func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Correct Request",
			req: createReqBodyToRequest(CreateRequestBody{
				ID:  uuid.New(),
				URL: "https://example.com",
				StreamState: room.StreamState{
					CurrentTime:  314.159,
					Paused:       true,
					PlaybackRate: 1.0,
				},
				StreamElement: nil,
				Timestamp:     time.Now().UTC(),
			}),
			_func: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
				}
				resBody := CreateResponseBody{Code: ""}
				err := json.NewDecoder(rec.Body).Decode(&resBody)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}
			},
		},
		{
			name: "Empty Request",
			req:  httptest.NewRequest(http.MethodPost, "/create", nil),
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
			CreateHandler(rec, test.req)
			test._func(t, rec)
		})
	}
}
