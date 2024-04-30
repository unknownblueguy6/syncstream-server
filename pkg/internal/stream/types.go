package stream

type StreamState struct {
	CurrentTime  float64 `json:"currentTime"  validate:"numeric"`
	Paused       bool    `json:"paused" validate:"boolean"`
	PlaybackRate float32 `json:"playbackRate" validate:"required,numeric"`
}

type StreamElement any
