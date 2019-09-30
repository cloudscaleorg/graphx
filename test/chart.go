package test

import (
	"fmt"

	"github.com/cloudscaleorg/graphx"
)

func GetCharts(n int) []*graphx.Chart {
	out := []*graphx.Chart{}
	for i := 0; i < n; i++ {
		out = append(out, &graphx.Chart{
			Name: fmt.Sprintf("test-chart-%d", i),
			Metrics: []graphx.ChartMetric{
				graphx.ChartMetric{
					Name:       fmt.Sprintf("test-metric-%d", i),
					Chart:      fmt.Sprintf("test-chart-%d", i),
					Query:      fmt.Sprintf("test-query-%d", i),
					DataSource: fmt.Sprintf("test-datasource-%d", i),
				},
			},
		})
	}
	return out
}
