package metric

import (
	"fmt"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

const EmptyGatewayIP = "None"

type Generator struct {
	api       gardena.API
	gatewayIP string
}

// NewGenerator creates a new Generator with a given gardena.API and a gatewayIP as string
func NewGenerator(api gardena.API, gatewayIP string) *Generator {
	var generator Generator
	generator.api = api
	generator.gatewayIP = gatewayIP
	return &generator
}

// InitializeLocationsMetrics queries all locations and sets up a metric to expose the number of locations
func (generator *Generator) InitializeLocationsMetrics() error {
	locations, err := generator.api.GetLocations()
	if err != nil {
		return fmt.Errorf("unable to get locations, got errer: %w", err)
	}
	locationsTotal.WithLabelValues(generator.api.GetBaseURL()).Set(float64(len(locations.Data)))
	return nil
}

// MonitorHealthOfEndpoints checks if the configured api health endpoint and the gateway bridge device
// are healthy by querying the endpoint urls. The result is exported as metric.
// If no ip for the bridge device is configured, this endpoint is ignored.
func (generator *Generator) MonitorHealthOfEndpoints() {
	timer := prometheus.NewTimer(endpointHealthCheckDuration.WithLabelValues())
	defer timer.ObserveDuration()

	gardenaApiUp := 0
	gardenaApiHealthUrl := generator.api.GetAPIHealthURL()
	if up := checkHealth(gardenaApiHealthUrl); up {
		gardenaApiUp = 1
	}
	hostHealth.WithLabelValues("api", gardenaApiHealthUrl).Set(float64(gardenaApiUp))

	if generator.gatewayIP != EmptyGatewayIP {
		gardenaGatewayUp := 0
		gatewayUrl := "http://" + generator.gatewayIP
		if up := checkHealth(gatewayUrl); up {
			gardenaGatewayUp = 1
		}
		hostHealth.WithLabelValues("gateway", gatewayUrl).Set(float64(gardenaGatewayUp))
	}
}

// checkHealth performs an HTTP GET request against a given endpoint url. If the request fails or
// the status code isn't '200' the endpoint is considered unhealthy.
func checkHealth(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Received error querying endpoint '%s'! Err was '%v'\n", url, err)
		return false
	}
	if resp.StatusCode != 200 {
		log.Printf("Received status code '%v' of endpoint '%s'!", resp.StatusCode, url)
		return false
	}
	return true
}
