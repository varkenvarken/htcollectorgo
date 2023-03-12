package collector

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Reading struct {
	StationId   string
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}

type Station struct {
	StationId string `form:"stationid"`
	Name      string `form:"name"`
}

type DB struct {
	db  *sql.DB
	mtx sync.Mutex
}

// SaveReading saves the temperature, humidity and stationid for a reading to a database.
// It accepts three arguments, temperature and humidity as float64 and stationid as a string.
// It returns an error if it is not able to save the reading to the database.
// It is save to run in concurrent go routines
func (db *DB) SaveReading(temperature float64, humidity float64, stationid string) error {
	reading := Reading{
		StationId:   stationid,
		Temperature: temperature,
		Humidity:    humidity,
		Timestamp:   time.Now(),
	}

	stmt, err := db.db.Prepare("INSERT INTO readings(stationid, temperature, humidity, timestamp) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	db.mtx.Lock()
	defer db.mtx.Unlock()

	_, err = stmt.Exec(reading.StationId, reading.Temperature, reading.Humidity, reading.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

// GetReadingsSince retrieves all the readings for a given station since a certain time.
// It accepts two arguments, since as a time.Time and stationid as a string.
// It returns a slice of Reading structs and an error if it fails to retrieve the readings.
func (db *DB) GetReadingsSince(since time.Time, stationid string) ([]Reading, error) {
	rows, err := db.db.Query("SELECT stationid, timestamp, temperature, humidity FROM readings WHERE timestamp >= ? AND stationid = ?", since, stationid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []Reading

	for rows.Next() {
		var reading Reading
		err := rows.Scan(&reading.StationId, &reading.Timestamp, &reading.Temperature, &reading.Humidity)
		if err != nil {
			return nil, err
		}
		readings = append(readings, reading)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return readings, nil
}

// GetDistinctStations returns all distinct station ids as a slice of strings
func (db *DB) GetDistinctStations() ([]string, error) {
	rows, err := db.db.Query("SELECT DISTINCT(stationid) FROM readings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stationids []string

	for rows.Next() {
		var stationid string
		err := rows.Scan(&stationid)
		if err != nil {
			return nil, err
		}
		stationids = append(stationids, stationid)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stationids, nil
}

// GetStationNames returns a map of stationids to names for all stationids
// for which there are readings.
// if no name is associated with a stationid, the name is listed as "Unknown"
func (db *DB) GetStationNames() (map[string]string, error) {

	stationids, err := db.GetDistinctStations()
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query("SELECT stationid, name FROM stationidtoname")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stationidtoname = make(map[string]string)

	for rows.Next() {
		var stationid, name string
		err := rows.Scan(&stationid, &name)
		if err != nil {
			return nil, err
		}
		stationidtoname[stationid] = name
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, stationid := range stationids {
		_, ok := stationidtoname[stationid]
		if !ok {
			stationidtoname[stationid] = "Unknown"
		}
	}

	return stationidtoname, nil
}

func (db *DB) GetAllReadingsSince(since time.Time) (map[string][]Reading, error) {

	stationids, err := db.GetDistinctStations()
	if err != nil {
		return nil, err
	}

	stationreadings := make(map[string][]Reading)

	for _, stationid := range stationids {
		rows, err := db.db.Query("SELECT stationid, timestamp, temperature, humidity FROM readings WHERE timestamp >= ? AND stationid = ?", since, stationid)
		if err != nil {
			return nil, err
		}

		var readings []Reading

		for rows.Next() {
			var reading Reading
			err := rows.Scan(&reading.StationId, &reading.Timestamp, &reading.Temperature, &reading.Humidity)
			if err != nil {
				return nil, err
			}
			readings = append(readings, reading)
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
		rows.Close()
		stationreadings[stationid] = readings
	}
	return stationreadings, nil
}

// UpdateName inserts or updates a Station struct
func (db *DB) UpdateName(station *Station) error {

	stmt, err := db.db.Prepare("REPLACE INTO stationidtoname(stationid, name) VALUES(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	db.mtx.Lock()
	defer db.mtx.Unlock()

	_, err = stmt.Exec(station.StationId, station.Name)
	if err != nil {
		return err
	}

	return nil
}

func MustInitDB(dataSourceName string) *DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS readings(id INTEGER PRIMARY KEY, stationid TEXT, temperature REAL, humidity REAL, timestamp TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS ts ON readings(timestamp)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS ts ON readings(stationid)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS stationidtoname(stationid TEXT NOT NULL PRIMARY KEY, name TEXT NOT NULL)")
	if err != nil {
		log.Fatal(err)
	}

	return &DB{db: db}
}

func (db *DB) Close() error {
	return db.db.Close()
}
