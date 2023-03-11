package collector

import (
	"fmt"
	"net/http"
	"strconv"
    "encoding/json"
    "time"
)

func HandleStoreReading(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tempStr := r.URL.Query().Get("temp")
		humStr := r.URL.Query().Get("hum")
		if tempStr == "" || humStr == "" {
			http.Error(w, "missing temp or hum query parameter", http.StatusBadRequest)
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
		db.SaveReading(temp, hum)
		fmt.Fprintf(w, "stored reading: temp=%.1f hum=%.1f\n", temp, hum)
		Temperature.Set(temp)
		Humidity.Set(hum)
	}
}


func HandleGetReadings(db *DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Calculate the timestamp 24 hours ago
        since := time.Now().Add(-24 * time.Hour)

        // Retrieve readings from the database since the specified time
        readings, err := db.GetReadingsSince(since)
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
