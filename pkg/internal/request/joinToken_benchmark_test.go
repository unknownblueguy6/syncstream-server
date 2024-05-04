package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"syncstream-server/pkg/internal/room"
)

// BenchmarkJoinTokenEndpoint benchmarks the JoinTokenHandler for performance and efficiency.
func BenchmarkJoinTokenEndpoint(b *testing.B) {
	// Setup necessary environment for the test
	ss := room.StreamState{
		CurrentTime:  314.159,
		Paused:       true,
		PlaybackRate: 1.0,
	}
	code, err := room.Manager.AddRoom(uuid.New(), "http://example.com", ss, nil)
	if err != nil {
		b.Fatalf("Failed to add room: %v", err)
	}

	// Prepare the request body
	reqBody := JoinTokenRequestBody{
		ID:   uuid.New(), // Use a new UUID for each test to simulate different user requests
		Code: code,       // Use the valid code from the setup
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		b.Fatalf("Failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/joinToken", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Benchmark the JoinTokenHandler
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		JoinTokenHandler(rec, req)
		// Reset the Body of the request for the next iteration to ensure it's not consumed
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}
