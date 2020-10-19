package main

import "github.com/prometheus/client_golang/prometheus"

const namespace = "speedtest"

type SpeedTestCollector struct {
	st *SpeedTest

	download *prometheus.Desc
	upload   *prometheus.Desc
	ping     *prometheus.Desc
}

func NewSpeedTestCollector(st *SpeedTest) *SpeedTestCollector {
	download := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download"),
		"Download bandwidth (Mbps).",
		nil, nil,
	)
	upload := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload"),
		"Upload bandwidth (Mbps).",
		nil, nil,
	)
	ping := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ping"),
		"Latency (ms)",
		nil, nil,
	)

	return &SpeedTestCollector{
		st:       st,
		download: download,
		upload:   upload,
		ping:     ping,
	}
}

func (stc *SpeedTestCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- stc.download
	ch <- stc.upload
	ch <- stc.ping
}

func (stc *SpeedTestCollector) Collect(ch chan<- prometheus.Metric) {
	metrics := stc.st.Run()
	ch <- prometheus.MustNewConstMetric(stc.download, prometheus.GaugeValue, metrics.Download/1e6)
	ch <- prometheus.MustNewConstMetric(stc.upload, prometheus.GaugeValue, metrics.Upload/1e6)
	ch <- prometheus.MustNewConstMetric(stc.ping, prometheus.GaugeValue, metrics.Ping)
}
