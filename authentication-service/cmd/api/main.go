package main

import (
	"authentication-service/data"

	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the authentication service")

	// Connect to the database
	conn := connectToDB()
	if conn == nil {
		log.Println("Could not connect to the postgresql database")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Println("Error starting the server", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Error connecting to the database", err)
			counts++
		} else {
			log.Println("Connected to the postgres database")
			return connection
		}

		if counts > 10 {
			log.Println("Could not connect to the database after 10 attempts")
			return nil
		}

		log.Println("Trying to connect to the database again after two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
