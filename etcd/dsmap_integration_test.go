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

func Test_DataSourceStore_Integration_Store(t *testing.T) {
	table := []struct {
		name string
		ds   int
	}{
		{
			name: "1 datasource",
			ds:   1,
		},
		{
			name: "10 datasource",
			ds:   10,
		},
		{
			name: "50 datasource",
			ds:   50,
		},
		{
			name: "100 datasource",
			ds:   100,
		},
		{
			name: "500 datasource",
			ds:   500,
		},
		{
			name: "1000 datasource",
			ds:   1000,
		},
		{
			name: "5000 datasource",
			ds:   5000,
		},
	}
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			client, teardown := et.Setup(t, nil)
			defer teardown()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

			m, err := NewDSMap(ctx, client)
			assert.NoError(t, err)

			sources := test.GenDataSources(tt.ds)
			m.Store(sources)
			assert.NoError(t, err)

			// allow convergence
			time.Sleep(1 * time.Second)
			cancel()

			// confirm stored datastores
			m.mu.Lock()
			for _, ds := range sources {
				v, ok := m.m[ds.Name]
				if !ok {
					t.Fatalf("failed to find datasource: %v. dump: %v", ds.Name, m.m)
				}
				assert.Equal(t, ds, v)
			}
			m.mu.Unlock()
		})
	}

}
