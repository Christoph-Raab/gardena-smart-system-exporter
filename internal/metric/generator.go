package metric

import (
	"fmt"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/internal/state"
	"github.com/Christoph-Raab/gardena-smart-system-exporter/pkg/gardena"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

const EmptyGatewayIP = "None"

type Generator struct {
	api       gardena.API
	store     state.Store
	gatewayIP string
}

// NewGenerator creates a new Generator with a given gardena.API and a gatewayIP as string
func NewGenerator(api gardena.API, gatewayIP string) *Generator {
	var g Generator
	g.api = api
	g.gatewayIP = gatewayIP
	g.store = state.NewStore()
	return &g
}

// InitializeLocationsMetrics queries all locations and for each location it adds the location's
// devices to the generator's store. It also sets up a metric about the number of locations.
func (g *Generator) InitializeLocationsMetrics() error {
	locations, err := g.api.GetLocations()
	if err != nil {
		return fmt.Errorf("unable to get locations, got errer:\n%w", err)
	}
	locationsTotal.WithLabelValues(g.api.GetBaseURL()).Set(float64(len(locations.Data)))

	for _, l := range locations.Data {
		// Returns: Ref test/location.json
		s, err := g.api.GetInitialStateFor(l.Location)
		if err != nil {
			return fmt.Errorf("getting initial state for location %s failed, got err:\n%w", s.Data.Id, err)
		}

		// list 6 objs (2 DEVICE, 2 COMMON, MOWER, SENSOR) -> store as 2 devices
		err = g.store.StoreDevices(*s)
		if err != nil {
			return fmt.Errorf("storing devices for location %s failed with err:\n%w", l.Id, err)
		}
	}
	return nil
}

// MonitorHealthOfEndpoints checks if the configured api health endpoint and the gateway bridge device
// are healthy by querying the endpoint urls. The result is exported as metric.
// If no ip for the bridge device is configured, this endpoint is ignored.
func (g *Generator) MonitorHealthOfEndpoints() {
	timer := prometheus.NewTimer(endpointHealthCheckDuration.WithLabelValues())
	defer timer.ObserveDuration()

	gardenaApiUp := 0
	gardenaApiHealthUrl := g.api.GetAPIHealthURL()
	if up := checkHealth(gardenaApiHealthUrl); up {
		gardenaApiUp = 1
	}
	hostHealth.WithLabelValues("api", gardenaApiHealthUrl).Set(float64(gardenaApiUp))

	if g.gatewayIP != EmptyGatewayIP {
		gardenaGatewayUp := 0
		gatewayUrl := "http://" + g.gatewayIP
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
