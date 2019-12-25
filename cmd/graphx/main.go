package main

import (
	"context"
	h "net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cloudscaleorg/graphx/admin"
	"github.com/cloudscaleorg/graphx/etcd"
	"github.com/cloudscaleorg/graphx/http"
	"github.com/cloudscaleorg/graphx/prometheus"
	"github.com/cloudscaleorg/graphx/registry"
	"github.com/crgimenes/goconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	v3 "go.etcd.io/etcd/clientv3"
)

// Config this struct is using the goconfig library for simple flag and env var
// parsing. See: https://github.com/crgimenes/goconfig
type Config struct {
	AdminListenAddr     string `cfg:"admin" cfgHelper:"the address in host:port format where the admin Api will listen. optional" cfgDefault:"0.0.0.0:8080"`
	WebsocketListenAddr string `cfg:"websocket" cfgHelper:"the address in host:port format where the websocket API will listen." cfgDefault:"0.0.0.0:8081"`
	Etcd                string `cfg:"etcd" cfgHelper:"a comma separated list of etcd hosts in host:port format. required" cfgRequired:"true" cfgDefault:"localhost:2379"`
	LogLevel            string `cfg:"debug" cfgHelper:"the debug level to use" cfgDefault:"debug"`
}

func main() {
	goconfig.PrefixEnv = "graphx"
	var conf Config
	err := goconfig.Parse(&conf)
	if err != nil {
		log.Fatal().Msgf("failed to parse any configuration: %v", err)
	}
	zerolog.SetGlobalLevel(logLevel(conf))
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	beReg := beRegInit(context.TODO())
	a, err := adminInit(context.TODO(), conf.Etcd, beReg)
	if err != nil {
		log.Fatal().Msgf("failed to create admin interface: %v", err)
	}
	adminServer := adminServer(conf.AdminListenAddr, a)
	if err != nil {
		log.Fatal().Msgf("failed to create admin api: %v", err)
	}
	eC := make(chan error, 2)
	sigs := signalHandler()
	go func() {
		log.Info().Msgf("starting admin http server on: %v", conf.AdminListenAddr)
		err := adminServer.ListenAndServe()
		if err != nil {
			eC <- err
		}
	}()
	select {
	case e := <-eC:
		log.Fatal().Msgf("received error: %v", e)
	case s := <-sigs:
		log.Info().Msgf("received signal: %v. stopping", s)
		adminServer.Shutdown(context.TODO())
	}
}

// creates an http server with the endpoints for administering a graphx cluster
func adminServer(addr string, admin admin.All) *h.Server {
	mux := h.NewServeMux()
	mux.HandleFunc("/api/v1/charts", http.ChartCRUD(admin))
	mux.HandleFunc("/api/v1/datasources", http.DataSourceCRUD(admin))
	mux.HandleFunc("/api/v1/backends", http.BackendCRUD(admin))
	s := &h.Server{
		Addr:    addr,
		Handler: mux,
	}
	return s
}

// creates a backend registry and registers implemented backends
func beRegInit(ctx context.Context) registry.Backend {
	beReg := registry.NewBackendReg()
	beReg.Register(prometheus.Prometheus, prometheus.NewBackend)
	return beReg
}

// creates the etcd backed admin interface
func adminInit(ctx context.Context, hosts string, beReg registry.Backend) (admin.All, error) {
	endpoints := strings.Split(hosts, ",")
	client, err := v3.New(
		v3.Config{
			Endpoints: endpoints,
		},
	)
	if err != nil {
		return nil, err
	}
	dMap, err := etcd.NewDSMap(ctx, client)
	if err != nil {
		return nil, err
	}
	cMap, err := etcd.NewChartMap(ctx, client)
	if err != nil {
		return nil, err
	}
	a := admin.NewAdmin(dMap, cMap, beReg)
	return a, nil
}

func signalHandler() <-chan os.Signal {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	return sigs
}

func logLevel(conf Config) zerolog.Level {
	level := strings.ToLower(conf.LogLevel)
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
