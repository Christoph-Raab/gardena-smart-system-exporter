package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

const metricNameSpace = "gardena_smart_system"
const DefaultGatewayIP = "None"

func Generate(gatewayIP string) bool {
	timer := prometheus.NewTimer(scrapeDuration.WithLabelValues())
	defer timer.ObserveDuration()

	gardenaApiUp := 0
	gardenaApiHealthUrl := "https://api.smart.gardena.dev/v1/health"
	if up := checkHealth(gardenaApiHealthUrl); up {
		gardenaApiUp = 1
	}
	hostHealth.WithLabelValues("api", gardenaApiHealthUrl).Set(float64(gardenaApiUp))

	if gatewayIP != DefaultGatewayIP {
		gardenaGatewayUp := 0
		gatewayUrl := "http://" + gatewayIP
		if up := checkHealth(gatewayUrl); up {
			gardenaGatewayUp = 1
		}
		hostHealth.WithLabelValues("gateway", gatewayUrl).Set(float64(gardenaGatewayUp))
	}

	return true
}

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
