package request

import (
	"log/slog"
	"syncstream-server/pkg/internal/room"
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func TestJoinHandler(t *testing.T) {
	ss := room.StreamState{
		CurrentTime:  314.159,
		Paused:       true,
		PlaybackRate: 1.0,
	}
	id := uuid.New()
	code, _ := room.Manager.AddRoom(id, "http://example.com", ss, nil)
	validToken := uuid.New()
	room.Manager.Tokens[validToken] = room.EphemeralTokenData{
		ID:         id,
		Code:       code,
		ExpiryTime: time.Now().Add(6 * time.Hour),
	}
	invalidToken := uuid.New()
	expiredToken := uuid.New()
	room.Manager.Tokens[expiredToken] = room.EphemeralTokenData{
		ID:         id,
		Code:       code,
		ExpiryTime: time.Now().Add(-6 * time.Hour),
	}

	go room.Manager.Run()
	server := httptest.NewServer(http.HandlerFunc(JoinHandler))
	slog.SetLogLoggerLevel(slog.LevelDebug)
	defer server.Close()

	tests := []struct {
		name  string
		token string
		code  int
	}{
		{
			name:  "Valid Token",
			token: validToken.String(),
			code:  http.StatusOK,
		},
		{
			name:  "Invalid Token",
			token: invalidToken.String(),
			code:  http.StatusForbidden,
		},
		{
			name:  "Expired Token",
			token: expiredToken.String(),
			code:  http.StatusForbidden,
		},
		{
			name:  "Empty Token",
			token: "12345",
			code:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testURL, _ := url.Parse(server.URL)
			testURL.Scheme = "ws"
			testURL = testURL.JoinPath("/join")
			q := testURL.Query()
			q.Add("token", test.token)
			testURL.RawQuery = q.Encode()
			_, res, err := websocket.DefaultDialer.Dial(testURL.String(), nil)
			if err != nil {
				if res.StatusCode != test.code {
					t.Errorf("expected status code %d, got %d", test.code, res.StatusCode)
				}
			}
		})
	}
}
