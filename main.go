package main

import (
	"log"
	"net/http"
	"swift-codes-project/db"
	handler "swift-codes-project/handlers"
	"swift-codes-project/parser"
	"swift-codes-project/service"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the SQLite database. The DSN specifies that the database is stored in a file.
	database, err := db.InitDB("file:swift_codes.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("Could not initialize DB: %v", err)
	}
	defer database.Close()

	// parse and store data from the XLSX file.
	err = parser.ParseExcelAndStore(database, "data/SWIFT_CODES.xlsx")
	if err != nil {
		log.Printf("Failed to parse/store Excel data: %v", err)
	}

	repo := &service.SwiftRepository{DB: database}
	httpHandler := &handler.SwiftHTTPHandler{DataStore: repo}
	// Set up the router for HTTP endpoints
	router := mux.NewRouter()
	router.HandleFunc("/v1/swift-codes/{code}", httpHandler.GetSwiftCode).Methods("GET")
	router.HandleFunc("/v1/swift-codes/country/{iso2}", httpHandler.GetCountrySwiftCodes).Methods("GET")
	router.HandleFunc("/v1/swift-codes", httpHandler.CreateSwiftCode).Methods("POST")
	router.HandleFunc("/v1/swift-codes/{code}", httpHandler.DeleteSwiftCode).Methods("DELETE")

	log.Println("Server is starting on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
