package graphx

// Streamer is an interface providing a streaming API to clients
type Streamer interface {
	// Recv blocks until either a metric or an error is available.
	Recv() (*Metric, error)
}
