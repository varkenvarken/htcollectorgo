package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	mux := http.NewServeMux()
	mux.HandleFunc("/storereading", collector.HandleStoreReading(db))
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/getreadings", collector.HandleGetReadings(db))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if server.ListenAndServe() != nil {
		log.Fatal().Str("Addr", server.Addr).Msg("Cannot listen on address. Terminating htcollector")
	}
}
