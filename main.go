package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TimeResponse struct {
	TorontoTime string `json:"toronto_time"`
}

func main() {
	http.HandleFunc("/current-time", timeHandler)
	http.HandleFunc("/display-time", Handler)

	http.ListenAndServe(":9999", nil)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	torontoTime := getCurrentTorontoTime()
	saveTimeToDatabase(torontoTime)

	response := TimeResponse{TorontoTime: torontoTime.Format(time.RFC3339)}
	json.NewEncoder(w).Encode(response)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	retrieveDatabase()
}

func getCurrentTorontoTime() time.Time {
	loc, _ := time.LoadLocation("America/Toronto")
	return time.Now().In(loc)
}

func saveTimeToDatabase(time time.Time) {
	db, err := sql.Open("mysql", "smithan:smithan@/goapi")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO time_log (time) VALUES (?)", time)
	if err != nil {
		panic(err)
	}
}

func retrieveDatabase() {
	// Database connection parameters
	db, err := sql.Open("mysql", "smithan:smithan@/goapi")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Fetch data from the time_table
	rows, err := db.Query("SELECT id, time FROM time_log")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var id int
		var timeValue string // use string to scan into

		err := rows.Scan(&id, &timeValue)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the time string into time.Time
		parsedTime, err := time.Parse("2006-01-02 15:04:05", timeValue)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %d, Time: %s\n", id, parsedTime.Format("2006-01-02 15:04:05"))
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
