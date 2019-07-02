//+build integration

package prometheus

import (
	"context"
	"testing"
	"time"

	promClient "github.com/prometheus/client_golang/api"
	promAPI "github.com/prometheus/client_golang/api/prometheus/v1"
	promModels "github.com/prometheus/common/model"
)

const (
	DefaultPromURL = "http://localhost:9090/"
)

func setupAPI(t *testing.T, url string) promAPI.API {
	// create api client
	conf := promClient.Config{
		Address: url,
	}
	client, err := promClient.NewClient(conf)
	if err != nil {
		t.Fatalf("failed to create prom api client: %v", err)
	}

	// create api bindings
	api := promAPI.NewAPI(client)
	return api
}

var TestQueryTT = []struct {
	// name of test in test table
	name string
	// url used to connect to api
	url string
	// query string used in our test query
	query string
}{
	// this queries prometheis itself so should always suceeed
	{
		name:  "prometheus self query",
		url:   DefaultPromURL,
		query: `prometheus_engine_queries_concurrent_max{job="prometheus"}`,
	},
	// this query assumes local dev is running telegraf and grafana containers
	{
		name:  "cpu percentage, two containers",
		url:   DefaultPromURL,
		query: `docker_container_cpu_usage_percent{container_name=~"metrics_telegraf_1|metrics_grafana_1|"}`,
	},
}

// TestPoll tests a poll of prometheus. This can be used for reference of how the client API works
func TestQuery(t *testing.T) {
	for _, tt := range TestQueryTT {

		t.Run(tt.name, func(t *testing.T) {
			api := setupAPI(t, tt.url)

			// attempt query
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)

			value, err := api.Query(ctx, tt.query, time.Now())
			if err != nil {
				cancel()
				t.Fatalf("failed to query API: %v", err)
			}

			vector, ok := value.(promModels.Vector)
			if !ok {
				cancel()
				t.Fatalf("failed to type cast returned value into Vector.")
			}
			if len(vector) <= 0 {
				t.Fatalf("retrieved zero metrics for query")
			}

			for _, sample := range vector {
				t.Logf("Metric: %v", sample.Metric)
				t.Logf("Value: %v", sample.Value)
				t.Logf("TimeStamp: %v", sample.Timestamp)
				t.Logf("Name: %s", sample.Metric["container_name"])
			}
			cancel()
		})

	}
}

var TestQueryRangeTT = []struct {
	// name of test in test table
	name string
	// url used to connect to api
	url string
	// query string used in our test query
	query string
}{
	// this queries prometheis itself so should always suceeed
	{
		name:  "prometheus self query",
		url:   DefaultPromURL,
		query: `prometheus_engine_queries_concurrent_max{job="prometheus"}[10s]`,
	},
	// this query assumes local dev is running telegraf and grafana containers
	{
		name:  "cpu percentage, two containers",
		url:   DefaultPromURL,
		query: `docker_container_cpu_usage_percent{container_name=~"metrics_telegraf_1|metrics_grafana_1|"}[10s]`,
	},
}

func TestQueryRange(t *testing.T) {
	for _, tt := range TestQueryRangeTT {

		t.Run(tt.name, func(t *testing.T) {
			api := setupAPI(t, tt.url)

			// attempt query
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)

			value, err := api.Query(ctx, tt.query, time.Now())
			if err != nil {
				cancel()
				t.Fatalf("failed to query API: %v", err)
			}

			matrix, ok := value.(promModels.Matrix)
			if !ok {
				cancel()
				t.Fatalf("failed to type cast returned value into Vector.")
			}
			if len(matrix) <= 0 {
				t.Fatalf("retrieved zero metrics for query")
			}

			cancel()
		})

	}
}
