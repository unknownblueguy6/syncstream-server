package stream

type StreamState struct {
	CurrentTime  float64 `json:"currentTime"`
	Paused       bool    `json:"paused"`
	PlaybackRate float32 `json:"playbackRate"`
}

type StreamElement any
