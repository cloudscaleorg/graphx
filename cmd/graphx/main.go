package main

import (
	"log"
	h "net/http"

	"github.com/cloudscaleorg/graphx/admin"
	"github.com/cloudscaleorg/graphx/http"
	"github.com/crgimenes/goconfig"
	v3 "go.etcd.io/etcd/clientv3"
)

// Config this struct is using the goconfig library for simple flag and env var
// parsing. See: https://github.com/crgimenes/goconfig
type Config struct {
	AdminListenAddr     string `cfg:"admin" cfgHelper:"the address in host:port format where the admin Api will listen. optional"`
	WebsocketListenAddr string `cfg:"websocket" cfgHelper:"the address in host:port format where the websocket API will listen. required" cfgRequired:"true"`
	Etcd                string `cfg:"etcd" cfgHelper:"a comma separated list of etcd hosts in host:port format. required" cfgRequired:"true"`
}

func main() {
	goconfig.PrefixEnv = "graphx"

	var conf Config
	err := goconfig.Parse(&conf)
	if err != nil {
		log.Fatalf("failed to parse any configuration: %v", err)
	}

}

func adminAPI(addr string, cAdmin admin.Chart, dAdmin admin.DataSource) *h.Server {
	mux := h.NewServeMux()
	mux.HandleFunc("/api/v1/chart", http.ChartCRUD(cAdmin))
	mux.HandleFunc("/api/v1/datasource", http.ChartCRUD(dAdmin))

	s := &h.Server{
		Addr:    addr,
		Handler: mux,
	}

	return h
}

func etcd(hosts string) *v3.Client {

}
