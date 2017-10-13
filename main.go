package main

import (
	"github.com/joho/godotenv"
	"database/sql"
	"os"
	"log"
	"time"
	"strconv"
)

func main() {
	// Load environment
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get MYSQL handle
	db, err := sql.Open("mysql", os.Getenv("SQL_USER") + ":" + os.Getenv("SQL_PASS") + "@tcp("+ os.Getenv("SQL_HOST") + ":" + os.Getenv("SQL_PORT")+")/"+os.Getenv("SQL_DB"))
	if err != nil {
		log.Fatal("Error making a connection to MySQL")
	}
	defer db.Close()

	maxConnections, err := strconv.Atoi(os.Getenv("SQL_MAX_CONNECTIONS"))
	db.SetMaxOpenConns(maxConnections)

	// Run task loop
	for {
		CheckTasks(db)
		time.Sleep(time.Millisecond * 500)
	}

}