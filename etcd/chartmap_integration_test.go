//+build etcdintegration

package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/cloudscaleorg/graphx/test"
	et "github.com/ldelossa/goframework/test/etcd"
	"github.com/stretchr/testify/assert"
)

func Test_ChartStore_Integration_Store(t *testing.T) {
	table := []struct {
		name string
		c    int
	}{
		{
			name: "1 chart",
			c:    1,
		},
		{
			name: "10 chart",
			c:    10,
		},
		{
			name: "50 chart",
			c:    50,
		},
		{
			name: "100 chart",
			c:    100,
		},
		{
			name: "500 chart",
			c:    500,
		},
		{
			name: "1000 chart",
			c:    1000,
		},
		{
			name: "5000 chart",
			c:    5000,
		},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			client, teardown := et.Setup(t, nil)
			defer teardown()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

			m, err := NewChartMap(ctx, client)
			assert.NoError(t, err)

			charts := test.GetCharts(tt.c)
			m.Store(charts)
			assert.NoError(t, err)

			// allow convergence
			time.Sleep(1 * time.Second)
			cancel()

			// confirm stored datastores
			m.mu.Lock()
			for _, ds := range charts {
				v, ok := m.m[ds.Name]
				if !ok {
					t.Fatalf("failed to find chart: %v. dump: %v", ds.Name, m.m)
				}
				assert.Equal(t, ds, v)
			}
			m.mu.Unlock()
		})
	}

}
