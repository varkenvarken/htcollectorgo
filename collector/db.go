package collector

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Reading struct {
	Temperature float64
	Humidity    float64
	Timestamp   time.Time
}

type DB struct {
	db  *sql.DB
	mtx sync.Mutex
}

func (db *DB) SaveReading(temperature float64, humidity float64) error {
	reading := Reading{
		Temperature: temperature,
		Humidity:    humidity,
		Timestamp:   time.Now(),
	}

	stmt, err := db.db.Prepare("INSERT INTO readings(temperature, humidity, timestamp) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	db.mtx.Lock()
	defer db.mtx.Unlock()

	_, err = stmt.Exec(reading.Temperature, reading.Humidity, reading.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetReadingsSince(since time.Time) ([]Reading, error) {
	rows, err := db.db.Query("SELECT timestamp, temperature, humidity FROM readings WHERE timestamp >= ?", since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []Reading

	for rows.Next() {
		var reading Reading
		err := rows.Scan(&reading.Timestamp, &reading.Temperature, &reading.Humidity)
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

func MustInitDB(dataSourceName string) *DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS readings(id INTEGER PRIMARY KEY, temperature REAL, humidity REAL, timestamp TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	return &DB{db: db}
}

func (db *DB) Close() error {
	return db.db.Close()
}
