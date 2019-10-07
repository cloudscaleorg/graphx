package graphx

import (
	"encoding/json"
)

// Chart is a user defined container of ChartMetrics. Provides a
// high level name for a chart and a container for all of a chart's metrics
type Chart struct {
	// a name for this chart. must be unique to the system
	Name string `json:"name"`
	// a map of datasource names by chart metrics
	DataSources map[string][]ChartMetric `json:data_sources`
	// // a list of chart metrics this high level chart comprises
	// Metrics []ChartMetric `json:"metrics"`
}

func (c *Chart) ToJSON() ([]byte, error) {
	b, err := json.Marshal(c)
	return b, err
}

func (c *Chart) FromJSON(b []byte) error {
	err := json.Unmarshal(b, c)
	return err
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
	DataSource string `json:"datasource"`
}
