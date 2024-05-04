package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkInitEndpoint(b *testing.B) {
	// Using a generic map to create an empty JSON object
	reqBody := map[string]interface{}{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/init", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	b.ResetTimer() // Reset the timer to ignore the setup code.
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		InitHandler(rec, req)
		// Reset the Body of the request for the next iteration to ensure it's not consumed
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
}
