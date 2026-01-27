package main

import (
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"

	"log"
)

var db *sqlx.DB
var config Config

func init() {
	var err error
	config, err = readConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	log.Println("Read config file successfully")

	schemaCache = buildSchemaCache(config.Collections)

	db, err = connectToDatabase()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connected successfully")

	err = migrateDatabase(db, config.Collections)
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	log.Println("Database migrated successfully")

}

func main() {
	defer db.Close()

	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("POST /{collection}", insertDocumentHandler)
	http.HandleFunc("POST /{collection}/", insertDocumentHandler)
	http.HandleFunc("GET /{collection}/{id}", getDocumentHandler)

	log.Printf("QuickStore Server starting on http://%v:%v", config.Host, config.Port)

	err := http.ListenAndServe(fmt.Sprintf("%v:%v", config.Host, config.Port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
