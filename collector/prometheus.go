package collector

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func PrometheusHandler() gin.HandlerFunc {
    h := promhttp.Handler()

    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}