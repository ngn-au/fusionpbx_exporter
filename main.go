package main

import (
	"net/http"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	user     = kingpin.Flag("user", "PostgreSQL username").Default("fusionpbx").String()
	password = kingpin.Flag("password", "PostgreSQL password").Default("password").String()
	dbname   = kingpin.Flag("dbname", "PostgreSQL database name").Default("fusionpbx").String()
	host     = kingpin.Flag("host", "PostgreSQL host").Default("localhost").String()
	port     = kingpin.Flag("port", "PostgreSQL port").Default("5432").String()
)

var reg = prometheus.NewPedanticRegistry()

func main() {
	kingpin.Parse()

	go func() {
		for {
			CollectMetrics()
			time.Sleep(10 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":8080", nil)
}
