package request

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/google/uuid"
	
)

//TODO: fix join token benchmark by creating dummy function for createid

func BenchmarkJoinTokenEndpoint(b *testing.B) {
    reqBody := JoinTokenRequestBody{
        ID:   uuid.New(),
        Code: "SampleCode", // Use a sample or mocked code valid for your logic
    }
    body, _ := json.Marshal(reqBody)
    req := httptest.NewRequest(http.MethodPost, "/joinToken", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        rec := httptest.NewRecorder()
        JoinTokenHandler(rec, req)
        // Reset the Body of the request for the next iteration to ensure it's not consumed
        req.Body = io.NopCloser(bytes.NewBuffer(body))
    }
}
