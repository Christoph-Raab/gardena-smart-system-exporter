package main

import (
	"flag"
	"fmt"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/internal/metric"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func main() {
	var gatewayIP string
	var metricInterval int
	var secretFilePath string
	flag.StringVar(&gatewayIP, "gateway-ip", metric.EmptyGatewayIP, "Ip of the Smart System Gateway Bridge Device, e.g. 192.168.178.24")
	flag.IntVar(&metricInterval, "metric-interval", 30, "Time between each metric generation run in seconds")
	flag.StringVar(&secretFilePath, "secret-file-path", "/etc/secrets/gardena-smart-system-exporter", "The path where client-id and client-secret files are stored.")
	flag.Parse()

	api, err := gardena.NewAPI().
		WithSecretFilePath(secretFilePath).
		Initialize()
	if err != nil {
		log.Fatalf("unable to initialize the api, got error:\n%v", err)
	}

	g := metric.NewGenerator(*api, gatewayIP)
	if err := g.InitializeLocationsMetrics(); err != nil {
		log.Fatalf("Unable to setup initial location metrics, got err:\n%v", err)
	}

	log.Println("Start serving metrics...")
	go func() {
		for {
			g.MonitorHealthOfEndpoints()
			time.Sleep(time.Duration(metricInterval) * time.Second)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Gardena Smart System Exporter</title></head>
			<body>
			<h1>Gardena Smart System Exporter</h1>
			<p><a href="https://github.com/Christoph-Raab/gardena-smart-system-exporter">View source code on GitHub</a></p>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 9093), nil))
}
