package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudscaleorg/graphx"
	v3 "go.etcd.io/etcd/clientv3"
)

const (
	cPrefix          = "/graphx/charts"
	cPrefixTemplate  = "/graphx/charts/%s"
	dsPrefix         = "/graphx/datasources"
	dsPrefixTemplate = "/graphx/datasources/%s"
	readyTO          = 5 * time.Second
)

func putCharts(ctx context.Context, etcd *v3.Client, charts []*graphx.Chart) error {
	for _, chart := range charts {
		b, err := chart.ToJSON()
		var (
			key = fmt.Sprintf(cPrefixTemplate, chart.Name)
			val = string(b)
		)
		_, err = etcd.Put(ctx, key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func putDS(ctx context.Context, etcd *v3.Client, sources []*graphx.DataSource) error {
	for _, source := range sources {
		b, err := source.ToJSON()
		var (
			key = fmt.Sprintf(dsPrefixTemplate, source.Name)
			val = string(b)
		)
		_, err = etcd.Put(ctx, key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func delCharts(ctx context.Context, etcd *v3.Client, names []string) error {
	for _, name := range names {
		key := fmt.Sprintf(cPrefixTemplate, name)
		_, err := etcd.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}

func delDS(ctx context.Context, etcd *v3.Client, names []string) error {
	for _, name := range names {
		key := fmt.Sprintf(dsPrefixTemplate, name)
		_, err := etcd.Delete(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}
