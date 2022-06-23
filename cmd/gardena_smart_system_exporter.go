package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

const metricNameSpace = "gardena_smart_system"

var (
	scrapeDuration = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: metricNameSpace,
		Name:      "gathering_duration",
		Help:      "The duration the gathering of all data took",
	}, []string{})
	hostHealth = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNameSpace,
		Name:      "host_health",
		Help:      "Indicates if a host is healthy",
	}, []string{
		"host",
		"addr",
	})
)

func checkHealth(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Received error querying endpoint '%s'!", url)
		return false
	}
	if resp.StatusCode != 200 {
		log.Printf("Received status code '%v' of endpoint '%s'!", resp.StatusCode, url)
		return false
	}
	return true
}

func createMetrics() bool {
	timer := prometheus.NewTimer(scrapeDuration.WithLabelValues())
	defer timer.ObserveDuration()

	gardenaApiUp := 0
	gardenaApiHealthUrl := "https://api.smart.gardena.dev/v1/health"
	if up := checkHealth(gardenaApiHealthUrl); up {
		gardenaApiUp = 1
	}

	gardenaGatewayUp := 0
	gatewayUrl := "http://192.168.178.24/"
	if up := checkHealth(gatewayUrl); up {
		gardenaGatewayUp = 1
	}

	hostHealth.WithLabelValues("api", gardenaApiHealthUrl).Set(float64(gardenaApiUp))
	hostHealth.WithLabelValues("gateway", gatewayUrl).Set(float64(gardenaGatewayUp))

	return true
}

func main() {
	log.Println("Start serving metrics...")
	go func() {
		for {
			if ok := createMetrics(); !ok {
				log.Println("Metric creation failed!")
			}
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 9093), nil))
}
