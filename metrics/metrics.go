package metrics

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const metricsPort = ":9100"

var Quantiles = map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}

func Expose() {
	log.Printf("Metrics exposed on %s", metricsPort)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(metricsPort, nil); err != nil {
		log.Fatalf("Error exposing metrics: %v", err)
	}
}
