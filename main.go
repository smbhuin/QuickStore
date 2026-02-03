package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const defaultConfigFile = "./config.json"
const defaultDatabaseFile = "./quickstore.db"

var configFile string
var databaseFile string

var rootMux *http.ServeMux

func init() {
	flag.StringVar(&configFile, "config", defaultConfigFile, "Path to config file")
	flag.StringVar(&databaseFile, "db", defaultDatabaseFile, "Path to database file")
	flag.Parse()

	var err error
	config, err = readConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	log.Println("Read config file successfully")

	schemaCache = buildSchemaCache(config.Collections)
	authCache = buildAuthCache(config)

	openapiSpec, err = buildOpenapiSpec(config)
	if err != nil {
		log.Fatalf("Error creating OpenAPI spec: %v", err)
	}

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
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("OPTIONS /{collection}", mockOptionsHandler)
	mux.HandleFunc("OPTIONS /{collection}/{id}", mockOptionsHandler)
	mux.HandleFunc("GET /{collection}", getAllDocumentsHandler)
	mux.HandleFunc("POST /{collection}", insertDocumentHandler)
	mux.HandleFunc("POST /{collection}/", insertDocumentHandler)
	mux.HandleFunc("GET /{collection}/{id}", getDocumentHandler)
	mux.HandleFunc("PUT /{collection}/{id}", replaceDocumentHandler)
	mux.HandleFunc("DELETE /{collection}/{id}", deleteDocumentHandler)

	apiMux := SetGlobalHeaders(mux)

	rootMux = http.NewServeMux()
	rootMux.Handle("/api/", http.StripPrefix("/api", apiMux))
	rootMux.Handle("/docs/", http.StripPrefix("/docs", SwaggerHandler()))

}

func main() {
	defer db.Close()

	registerRoutes()

	log.Printf("QuickStore Server starting on http://%s:%d", config.Host, config.Port)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), rootMux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
