package main

import (
	"net/http"
	"time"

	"github.com/MaibornWolff/maDocK8s/exporter/sample-metrics/swagger/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "my_custom_metrtic",
		Help: "Custom metrich which increases every 2s",
	})
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/greetings", api.GreetingHandler)

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":2112", handler)
}
