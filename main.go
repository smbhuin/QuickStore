package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const defaultConfigFile = "./config.json"
const defaultDatabaseFile = "./quickstore.db"

var rootMux *http.ServeMux

func initApp(configFile string, databaseFile string) error {
	var err error

	config, err = readConfig(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	schemaCache = buildSchemaCache(config.Collections)
	authCache = buildAuthCache(config)

	openapiSpec, err = buildOpenapiSpec(config)
	if err != nil {
		return fmt.Errorf("error creating OpenAPI spec: %w", err)
	}

	db, err = connectToDatabase(databaseFile)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	err = migrateDatabase(db, config.Collections)
	if err != nil {
		return fmt.Errorf("error migrating database: %w", err)
	}

	return nil
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
	var configFile string
	var databaseFile string

	flag.StringVar(&configFile, "config", defaultConfigFile, "Path to config file")
	flag.StringVar(&databaseFile, "db", defaultDatabaseFile, "Path to database file")
	flag.Parse()

	err := initApp(configFile, databaseFile)
	if err != nil {
		log.Fatalf("Startup error: %v", err)
	}

	defer db.Close()

	registerRoutes()

	log.Printf("QuickStore Server starting on http://%s:%d", config.Host, config.Port)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), rootMux)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
