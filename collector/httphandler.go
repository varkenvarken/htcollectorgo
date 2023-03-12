package collector

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus"
)

// Store a reading from a ShellyHT
// note that this does *not* follow gin's regular parameter syntax
func HandleStoreReading(c *gin.Context, db *DB) {
	tempStr := c.Query("temp")
	humStr := c.Query("hum")
	idStr := c.Query("id")
	if tempStr == "" || humStr == "" || idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing query parameter"})
		return
	}
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid temp query parameter"})
		return
	}
	hum, err := strconv.ParseFloat(humStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hum query parameter"})
		return
	}
	// TODO add check on length of stationid
	db.SaveReading(temp, hum, idStr)

	c.JSON(http.StatusOK, Reading{StationId: idStr, Temperature: temp, Humidity: hum})

	// prometheus gauges
	Temperature.With(prometheus.Labels{"stationid": idStr}).Set(temp)
	Humidity.With(prometheus.Labels{"stationid": idStr}).Set(hum)
}

func HandleReadings(c *gin.Context, db *DB) {
	// Calculate the timestamp 24 hours ago
	since := time.Now().Add(-24 * time.Hour)

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing stationid"})
		return
	}
	// Retrieve readings from the database since the specified time
	readings, err := db.GetReadingsSince(since, idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot retrieve readings"})
		return
	}

	c.JSON(http.StatusOK, readings)
}

func HandleAllReadings(c *gin.Context, db *DB) {
	// Calculate the timestamp 24 hours ago
	since := time.Now().Add(-24 * time.Hour)

	// Retrieve readings from the database since the specified time
	readings, err := db.GetAllReadingsSince(since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot retrieve readings"})
		return
	}

	c.JSON(http.StatusOK, readings)
}

func HandleAllNames(c *gin.Context, db *DB) {
	stationnames, err := db.GetStationNames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot retrieve stationnames"})
		return
	}

	c.JSON(http.StatusOK, stationnames)
}

func HandleName(c *gin.Context, db *DB) {
	var station Station

	if err := c.ShouldBind(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.UpdateName(&station); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update stationname"})
		return
	} else {
		c.JSON(http.StatusCreated, station)
	}
}
