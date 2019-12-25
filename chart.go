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
	DataSources map[string][]*ChartMetric `json:"data_sources"`
}

func (c *Chart) ToJSON() ([]byte, error) {
	b, err := json.Marshal(c)
	return b, err
}

func (c *Chart) FromJSON(b []byte) error {
	err := json.Unmarshal(b, c)
	return err
}

// MergeDataSources takes 0 - N charts and merges duplicate datasource:chartMetric maps
// into one.
func MergeCharts(charts []*Chart) (map[string][]*ChartMetric, []string) {
	metricMap := map[string][]*ChartMetric{}
	dsNames := []string{}
	for _, chart := range charts {
		for datasource, chartMetrics := range chart.DataSources {
			metricMap[datasource] = append(metricMap[datasource], chartMetrics...)
		}
	}
	for dsName, _ := range metricMap {
		dsNames = append(dsNames, dsName)
	}
	return metricMap, dsNames
}
