package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	//"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/varkenvarken/htcollectorgo/collector"
)

func setLogLevel() {
	loglevel := zerolog.InfoLevel
	switch strings.ToLower(os.Getenv("LOGLEVEL")) {
	case "debug":
		loglevel = zerolog.DebugLevel
	case "error":
		loglevel = zerolog.ErrorLevel
	default:
		loglevel = zerolog.InfoLevel
	}
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger().Level(loglevel)
}

func main() {

	setLogLevel()
	log.Info().Msg("Starting htcollector")

	db := collector.MustInitDB("readings.db")
	defer db.Close()

	r := gin.Default()

	r.GET("/storereading", func(c *gin.Context) {
		collector.HandleStoreReading(c, db)
	})

	r.GET("/readings/:id", func(c *gin.Context) {
		collector.HandleReadings(c, db)
	})

	r.GET("/readings", func(c *gin.Context) {
		collector.HandleAllReadings(c, db)
	})

	r.GET("/names", func(c *gin.Context) {
		collector.HandleAllNames(c, db)
	})

	r.POST("/names", func(c *gin.Context) {
		collector.HandleName(c, db)
	})

	r.GET("metrics", collector.PrometheusHandler())

	r.Run(":1883")
}
