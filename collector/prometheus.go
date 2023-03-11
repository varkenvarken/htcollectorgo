package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	Temperature = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "temperature",
		Help: "The current temperature reading",
	})

	Humidity = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "humidity",
		Help: "The current humidity reading",
	})
)

func init() {
	prometheus.MustRegister(Temperature)
	prometheus.MustRegister(Humidity)
}
