package graphx

// Chart is the user provided chart configuration.
// a ChartDescriptor will one or more Chart objects.
// each Chart defines how and where graphx aggregates metrics.
type Chart struct {
	// a name for this chart. must be unique to the system
	Name string `json:"name"`
	// a map assocatiating datasource names with a list of ChartMetric.
	ChartMetrics map[string][]ChartMetric `json:"metrics"`
}

// ChartMetric
type ChartMetric struct {
	// name of the metric within a chart
	Name string `json:"name"`
	// the name of the chart this metric is destined for
	Chart string `json:"chart"`
	// the query to retrieve this metric
	Query string `json:"query"`
	// the datasource that the query targets
	Datasource string `json:"datasource"`
}

// MergeChartMetrics takes a list of charts and merges any ChartMetrics of the same
// datasource. This is useful when passing ChartMetric arrays to Queriers.
func MergeChartMetrics(charts []Chart) map[string][]ChartMetric {
	res := map[string][]ChartMetric{}

	for _, chart := range charts {
		for datasource, cmetrics := range chart.ChartMetrics {
			res[datasource] = append(res[datasource], cmetrics...)
		}
	}

	return res
}
