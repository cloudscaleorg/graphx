//+build integration

package prometheus

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/cloudscaleorg/graphx"
)

var TestQueryStreamerTT = []struct {
	// name of the test table
	name string
	// id of unique query session
	id string
	// url of prometheus instance for creating API client
	url string
	// // the names we expect to see in returned metrics from Querier
	// Names []string
	// string we used to set the poll interface
	pollInterval string
	// // the charts we expect to see in returned metrics from Querier
	// charts []graphx.ChartName
	// chart descriptor
	chartDesc graphx.ChartsDescriptor
}{
	{
		name:         "single chart, single name",
		id:           "integration-test",
		url:          DefaultPromURL,
		pollInterval: "2s",
		chartDesc: graphx.ChartsDescriptor{
			ChartNames: []graphx.ChartName{graphx.CPUPercentage},
			Names:      []string{"metrics_prometheus_1"},
		},
	},
	{
		name:         "single chart, two names",
		id:           "integration-test",
		url:          DefaultPromURL,
		pollInterval: "2s",
		chartDesc: graphx.ChartsDescriptor{
			ChartNames: []graphx.ChartName{graphx.CPUPercentage},
			Names:      []string{"metrics_prometheus_1", "metrics_grafana_1"},
		},
	},
	{
		name:         "single chart, three names",
		id:           "integration-test",
		url:          DefaultPromURL,
		pollInterval: "2s",
		chartDesc: graphx.ChartsDescriptor{
			ChartNames: []graphx.ChartName{graphx.CPUPercentage},
			Names:      []string{"metrics_prometheus_1", "metrics_grafana_1", "metrics_telegraf_1"},
		},
	},
	{
		name:         "two charts, single name",
		id:           "integration-test",
		url:          DefaultPromURL,
		pollInterval: "2s",
		chartDesc: graphx.ChartsDescriptor{
			ChartNames: []graphx.ChartName{graphx.CPUPercentage, graphx.MemUsage},
			Names:      []string{"metrics_prometheus_1"},
		},
	},
	{
		name:         "two charts, two names",
		id:           "integration-test",
		url:          DefaultPromURL,
		pollInterval: "2s",
		chartDesc: graphx.ChartsDescriptor{
			ChartNames: []graphx.ChartName{graphx.CPUPercentage, graphx.MemUsage},
			Names:      []string{"metrics_prometheus_1", "metrics_grafana_1"},
		},
	},
	{
		name:         "two charts, three names",
		id:           "integration-test",
		url:          DefaultPromURL,
		pollInterval: "2s",
		chartDesc: graphx.ChartsDescriptor{
			ChartNames: []graphx.ChartName{graphx.CPUPercentage, graphx.MemUsage},
			Names:      []string{"metrics_prometheus_1", "metrics_grafana_1", "metrics_telegraf_1"},
		},
	},
}

func confirmNamesSeen(t *testing.T, seenMap map[string]bool, Names []string) {
	// confirm we did not see an unexpected name
	for seenJH, _ := range seenMap {
		var expected bool
		for _, expectedJH := range Names {
			if seenJH == expectedJH {
				expected = true
			}
		}
		if !expected {
			t.Fatalf("saw unexpected name in seenMap: %v", seenMap)
		}
	}

	// confirm we saw the expeced names
	for _, Name := range Names {
		if _, ok := seenMap[Name]; !ok {
			t.Fatalf("expected to see name: %v but did not", Name)
		}
	}
	return
}

func confirmChartsSeen(t *testing.T, seenMap map[graphx.ChartName]bool, charts []graphx.ChartName) {
	// confirm we did not see an unexpected chart
	for seenChart, _ := range seenMap {
		var expected bool
		for _, expectedChart := range charts {
			if seenChart == expectedChart {
				expected = true
			}
		}
		if !expected {
			t.Fatalf("saw unexpected name in seenMap: %v", seenMap)
		}
	}

	// confirm we saw the expected charts
	for _, chart := range charts {
		if _, ok := seenMap[chart]; !ok {
			t.Fatalf("expected to see chart: %v but did not", chart)
		}
	}
	return
}

func TestQueryStreamer(t *testing.T) {
	for _, tt := range TestQueryStreamerTT {

		t.Run(tt.name, func(t *testing.T) {

			// create channel to view pass to querier
			mChan := make(chan *graphx.Metric, 1024)
			errChan := make(chan error, 1024)

			api := setupAPI(t, tt.url)

			// create querier
			fill := time.Now().Add(-1 * time.Hour)
			poll, err := time.ParseDuration(tt.pollInterval)
			if err != nil {
				t.Fatalf("failed to parse pollInterval into duration: %v", err)
			}

			// create query map
			dur, err := time.ParseDuration(tt.pollInterval)
			tt.chartDesc.PollInterval = graphx.Duration(dur)
			tt.chartDesc.Fill = graphx.TimeStamp(fill)
			tt.chartDesc.PollInterval = graphx.Duration(poll)

			qm, err := newQueryMap(tt.chartDesc)
			if err != nil {
				t.Fatalf("provided chart description failed to create query map: %v", err)
			}

			q := &queryStreamer{
				client:       api,
				cd:           tt.chartDesc,
				id:           tt.id,
				queryMap:     qm,
				mChan:        mChan,
				errChan:      errChan,
				lm:           &sync.RWMutex{},
				latestTS:     time.Now(),
				eosThreshold: 30 * time.Second,
			}

			// create context with acceptable timeout
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			// go q.Query(ctx)
			go q.fillBuf(ctx)

			NamesSeen := make(map[string]bool)
			chartsSeen := make(map[graphx.ChartName]bool)

			// receive off channel until error
			for {
				m, testErr := q.Recv(ctx)
				if testErr != nil {
					// fail if we see error other then ctx done
					if _, ok := testErr.(*graphx.CtxDoneErr); !ok {
						t.Fatalf("received unexpected error from stream: %v", testErr)
					}
					break
				}
				log.Printf("received metric: %v", m)
				NamesSeen[m.Name] = true
				chartsSeen[m.Chart] = true
			}

			confirmNamesSeen(t, NamesSeen, tt.chartDesc.Names)
			confirmChartsSeen(t, chartsSeen, tt.chartDesc.ChartNames)

			// cancel context before next test table
			cancel()
		})

	}
}
