package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "OK"}`)
}

func insertDocumentHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	w.Header().Set("Content-Type", "application/json")

	if !isCollectionExists(collectionName) {
		http.Error(w, `{"error": {"message": "Collection not found"}}`, http.StatusNotFound)
		return
	}

	var document map[string]any

	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&document)
	if err != nil {
		http.Error(w, `{"error": {"message": "Invalid JSON"}}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate against schema
	if !validateJSONByCollectionName(document, collectionName) {
		http.Error(w, `{"error": {"message": "Validation failed"}}`, http.StatusBadRequest)
		return
	}

	// Insert into database
	err = insertDocument(db, collectionName, document)
	if err != nil {
		log.Printf("Error inserting document: %v", err)
		http.Error(w, `{"error": {"message": "Failed to insert document"}}`, http.StatusInternalServerError)
		return
	}

	// Respond with success
	fmt.Fprintf(w, `{"message": "Document inserted"}`)
}

func getDocumentHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	idStr := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")
	if !isCollectionExists(collectionName) {
		http.Error(w, `{"error": {"message": "Collection not found"}}`, http.StatusNotFound)
		return
	}
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, `{"error": {"message": "Invalid ID"}}`, http.StatusBadRequest)
		return
	}
	document, err := getDocument(db, collectionName, id)
	if err != nil {
		log.Printf("Error retrieving document: %v", err)
		http.Error(w, `{"error": {"message": "Failed to retrieve document"}}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(document)
}
