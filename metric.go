package graphx

// A Metric is the datastructure we return to the client on request
type Metric struct {
	// Name is the name this particular metric is for
	Name string `json:"name"`
	// Chart is the chart this metric should be routed for for the given name
	Chart string `json:"chart_name"`
	// A timestamp in Unix format
	TimeStamp int64 `json:"time_stamp"`
	// the value to plot on the graph
	Value string `json:"value"`
}
