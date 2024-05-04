package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"syncstream-server/pkg/internal/room"
)

func BenchmarkCreateEndpoint(b *testing.B) {
	reqBody := CreateRequestBody{
		ID:  uuid.New(),
		URL: "https://example.com",
		StreamState: room.StreamState{
			CurrentTime:  314.159,
			Paused:       true,
			PlaybackRate: 1.0,
		},
		StreamElement: nil,
		Timestamp:     time.Now().UTC(),
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		CreateHandler(rec, req)
		// Reset the Body of the request for the next iteration to ensure it's not consumed
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}
