package collector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func HandleStoreReading(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tempStr := r.URL.Query().Get("temp")
		humStr := r.URL.Query().Get("hum")
		idStr := r.URL.Query().Get("id")
		if tempStr == "" || humStr == "" || idStr == "" {
			http.Error(w, "missing query parameter", http.StatusBadRequest)
			return
		}
		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			http.Error(w, "invalid temp query parameter", http.StatusBadRequest)
			return
		}
		hum, err := strconv.ParseFloat(humStr, 64)
		if err != nil {
			http.Error(w, "invalid hum query parameter", http.StatusBadRequest)
			return
		}
		// TODO add check on length of stationid
		db.SaveReading(temp, hum, idStr)
		fmt.Fprintf(w, "stored reading: id=%s temp=%.1f hum=%.1f\n", idStr, temp, hum)

		// prometheus gauges
		Temperature.With(prometheus.Labels{"stationid": idStr}).Set(temp)
		Humidity.With(prometheus.Labels{"stationid": idStr}).Set(hum)
	}
}

func HandleGetReadings(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Calculate the timestamp 24 hours ago
		since := time.Now().Add(-24 * time.Hour)

		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, "missing query parameter", http.StatusBadRequest)
			return
		}
		// Retrieve readings from the database since the specified time
		readings, err := db.GetReadingsSince(since, idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Encode the readings as JSON and write to the response
		jsonBytes, err := json.Marshal(readings)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)
	}
}

func HandleUpdateName(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stationid := r.URL.Query().Get("stationid")
		name := r.URL.Query().Get("name")
		if stationid == "" || name == ""{
			http.Error(w, "missing query parameter", http.StatusBadRequest)
			return
		}
		db.UpdateName(stationid, name)
		fmt.Fprintf(w, "updated: id=%s name=%s\n", stationid, name)
	}
}
