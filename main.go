package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

func main() {
	var err error
	db, err = sql.Open("sqlite", "file:db.sqlite3?_pragma=foreign_keys(1)&_time_format=datetime")
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	registerRoutes()

	log.Println("Abierto en http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", MiddlewareHeaders(http.DefaultServeMux)))
}
