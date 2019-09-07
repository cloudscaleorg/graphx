package test

import (
	"fmt"

	"github.com/cloudscaleorg/graphx"
)

func GenDataSources(n int) []*graphx.DataSource {
	out := []*graphx.DataSource{}
	for i := 0; i < n; i++ {
		out = append(out, &graphx.DataSource{
			Name:       fmt.Sprintf("test-datasource-%d", i),
			ConnString: "datasource://host:port?param=value",
		})
	}
	return out
}
