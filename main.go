package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cmd         = flag.String("cmd", "speedtest", "speedtest command")
	listenAddr  = flag.String("addr", ":9112", "address to listen on")
	metricsPath = flag.String("path", "/metrics", "metrics path")
)

func main() {
	flag.Parse()

	checkSpeedTestVersion()
	st := NewSpeedTest(*cmd, flag.Args())
	stc := NewSpeedTestCollector(st)

	prometheus.MustRegister(stc)
	http.Handle(*metricsPath, promhttp.Handler())

	log.Printf("listening on %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))

}
