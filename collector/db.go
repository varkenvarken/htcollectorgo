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

type DB struct {
	db  *sql.DB
	mtx sync.Mutex
}

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

func (db *DB) UpdateName(stationid, name string) error {

	stmt, err := db.db.Prepare("REPLACE stationidtoname(stationid, name) VALUES(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	db.mtx.Lock()
	defer db.mtx.Unlock()

	_, err = stmt.Exec(stationid, name)
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
