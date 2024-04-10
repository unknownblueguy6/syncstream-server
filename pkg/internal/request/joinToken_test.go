package request

import (
	"bytes"
	"encoding/json"
	"syncstream-server/pkg/internal/room"
	"syncstream-server/pkg/internal/stream"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

func joinTokenReqBodyToRequest(reqBody JoinTokenRequestBody) *http.Request {
	reqBodyBytes, _ := json.Marshal(reqBody)
	return httptest.NewRequest(http.MethodPost, "/joinToken", bytes.NewBuffer(reqBodyBytes))
}

func TestJoinTokenHandler(t *testing.T) {
	ss := stream.StreamState{
		CurrentTime:  314.159,
		Paused:       true,
		PlaybackRate: 1.0,
	}
	code, _ := room.Manager.AddRoom(uuid.New(), "http://example.com", ss, nil)
	tests := []struct {
		name  string
		req   *http.Request
		_func func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Correct Request",
			req: joinTokenReqBodyToRequest(JoinTokenRequestBody{
				ID:   uuid.New(),
				Code: code,
			}),
			_func: func(t *testing.T, rec *httptest.ResponseRecorder) {
				if rec.Code != http.StatusOK {
					t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
				}
				resBody := JoinTokenResponseBody{Token: uuid.UUID{}, ExpiryTime: time.Time{}}
				err := json.NewDecoder(rec.Body).Decode(&resBody)
				if err != nil {
					t.Errorf("failed to decode response body: %v", err)
				}
			},
		},
		{
			name: "Empty Request",
			req:  httptest.NewRequest(http.MethodPost, "/joinToken", nil),
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
			JoinTokenHandler(rec, test.req)
			test._func(t, rec)
		})
	}
}
