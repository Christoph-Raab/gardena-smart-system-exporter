package main

import (
	"flag"
	"fmt"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/internal/metric"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func main() {
	var gatewayIp string
	flag.StringVar(&gatewayIp, "gateway-ip", "None", "Ip of the Smart System Gateway, e.g. 192.168.178.24")
	flag.Parse()

	log.Println("Start serving metrics...")
	go func() {
		for {
			if ok := metric.Generate(gatewayIp); !ok {
				log.Println("Metric creation failed!")
			}
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 9093), nil))
}
