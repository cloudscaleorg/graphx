package graphx

import (
	"context"
	"time"
)

// StreamerFactory allows runtime creation of a Streamer.
// this is necessary in order to depedency inject a Streamer
// into the stream http handler.
type StreamerFactory interface {
	NewStreamer(ctx context.Context, id string, charts []*Chart, pollInterval time.Duration) Streamer
}
