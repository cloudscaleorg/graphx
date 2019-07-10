package graphx

// Chart is a user defined container of ChartMetrics. Provides a
// high level name for a chart and a container for all of a chart's metrics
type Chart struct {
	// a name for this chart. must be unique to the system
	Name string `json:"name"`
	// a list of chart metrics this high level chart comprises
	ChartMetrics []ChartMetric `json:"metrics"`
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

// DatasourceTranpose takes a list of charts and returns a map
// of ChartMetrics keye'd by their datasource. This is helpful for
// handing specific ChartMetrics to the appropriate datasource clients
func DatasourceTranspose(charts []*Chart) map[string][]ChartMetric {
	res := map[string][]ChartMetric{}

	for _, chart := range charts {
		for _, chartMetric := range chart.ChartMetrics {
			res[chartMetric.Datasource] = append(res[chartMetric.Datasource], chartMetric)
		}
	}

	return res
}
