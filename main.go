package main

import (
	"fmt"
	"log"
	"net/http"
)

const configFile = "config.json"
const databaseFile = "./quickstore.db"

func init() {
	var err error
	config, err = readConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	log.Println("Read config file successfully")

	schemaCache = buildSchemaCache(config.Collections)
	authCache = buildAuthCache(config)

	db, err = connectToDatabase(databaseFile)
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

func registerRoutes() {
	http.HandleFunc("GET /health", healthHandler)
	http.HandleFunc("GET /{collection}", getAllDocumentsHandler)
	http.HandleFunc("POST /{collection}", insertDocumentHandler)
	http.HandleFunc("POST /{collection}/", insertDocumentHandler)
	http.HandleFunc("GET /{collection}/{id}", getDocumentHandler)
	http.HandleFunc("PUT /{collection}/{id}", replaceDocumentHandler)
	http.HandleFunc("DELETE /{collection}/{id}", deleteDocumentHandler)
}

func main() {
	defer db.Close()

	registerRoutes()

	log.Printf("QuickStore Server starting on http://%s:%d", config.Host, config.Port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
