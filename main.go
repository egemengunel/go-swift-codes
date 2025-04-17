package main

import (
	"log"
	"net/http"
	"swift-codes-project/db"
	"swift-codes-project/parser"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the SQLite database. The DSN specifies that the database is stored in a file.
	database, err := db.InitDB("file:swift_codes.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Could not initialize DB: %v", err)
	}
	defer database.Close()

	// Optionally, parse and store data from the XLSX file.
	err = parser.ParseExcelAndStore(database, "data/SWIFT_CODES.xlsx")
	if err != nil {
		log.Printf("Failed to parse/store Excel data: %v", err)
	}

	// Set up the router for HTTP endpoints
	r := mux.NewRouter()
	// TODO: Register your API routes here

	log.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
