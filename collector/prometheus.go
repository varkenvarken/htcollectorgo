package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	Temperature = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "htcollector",
		Subsystem: "readings",
		Name:      "temperature",
		Help:      "The current temperature reading",
	}, []string{"stationid"})

	Humidity = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "htcollector",
		Subsystem: "readings",
		Name:      "humidity",
		Help:      "The current humidity reading",
	}, []string{"stationid"})
)
