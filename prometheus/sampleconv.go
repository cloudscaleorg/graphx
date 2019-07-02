package prometheus

import (
	"github.com/cloudscaleorg/graphx"
	promModels "github.com/prometheus/common/model"
)

const (
	// the metric tag used to extract our name. In our case containers are
	// created with our name name.
	NameTag = "container_name"
)

// samplePairToMetric converts a prometheus SamplePair to a graphx.Metric
func samplePairToMetric(chart string, metric promModels.Metric, sp promModels.SamplePair) *graphx.Metric {
	Name := string(metric[NameTag])
	m := &graphx.Metric{
		Chart:     chart,
		Name:      Name,
		TimeStamp: sp.Timestamp.Unix(),
		Value:     sp.Value.String(),
	}
	return m
}

// // samplePairToMetric converts a prometheus SamplePair to a graphx.Metric
// func samplePairToRangeMetric(sp promModels.SamplePair) *graphx.RangeMetric {
// 	m := &graphx.RangeMetric{
// 		TimeStamp: sp.Timestamp.Unix(),
// 		Value:     sp.Value.String(),
// 	}
// 	return m
// }

// sampleToMetric converts a prometheus Sample to our domain Metric object
func sampleToMetric(chart string, sample *promModels.Sample) *graphx.Metric {
	Name := string(sample.Metric[NameTag])
	m := &graphx.Metric{
		Chart:     chart,
		Name:      Name,
		TimeStamp: sample.Timestamp.Unix(),
		Value:     sample.Value.String(),
	}
	return m
}
