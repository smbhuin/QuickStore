package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getAuthTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}

	return ""
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, fmt.Sprintf(`{"error": {"message": "%s"}}`, message), code)
}

func sendSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "%s"}`, message)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	sendSuccess(w, "OK")
}

func insertDocumentHandler(w http.ResponseWriter, r *http.Request) {
	collectionName := r.PathValue("collection")
	if !isCollectionExists(collectionName) {
		sendError(w, "Collection not found", http.StatusNotFound)
		return
	}

	authToken := getAuthTokenFromRequest(r)
	if !isAuthTokenValid(authToken, collectionName, ActionCreate) {
		sendError(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	var document map[string]any
	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&document)
	if err != nil {
		sendError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validate against schema
	if !validateJSONByCollectionName(document, collectionName) {
		sendError(w, "Validation failed", http.StatusBadRequest)
		return
	}

	// Insert into database
	err = insertDocument(db, collectionName, document)
	if err != nil {
		log.Printf("Error inserting document: %v", err)
		sendError(w, "Failed to insert document", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, "Document inserted")
}

func getDocumentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collectionName := r.PathValue("collection")
	if !isCollectionExists(collectionName) {
		sendError(w, "Collection not found", http.StatusNotFound)
		return
	}

	authToken := getAuthTokenFromRequest(r)
	if !isAuthTokenValid(authToken, collectionName, ActionCreate) {
		http.Error(w, `{"error": {"message": "Unauthorized access"}}`, http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := StoiStrict(idStr)
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

func getAllDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collectionName := r.PathValue("collection")
	if !isCollectionExists(collectionName) {
		http.Error(w, `{"error": {"message": "Collection not found"}}`, http.StatusNotFound)
		return
	}

	authToken := getAuthTokenFromRequest(r)
	if !isAuthTokenValid(authToken, collectionName, ActionCreate) {
		http.Error(w, `{"error": {"message": "Unauthorized access"}}`, http.StatusUnauthorized)
		return
	}

	skip := Stoi(r.URL.Query().Get("skip"), 0)
	limit := rangeBound(Stoi(r.URL.Query().Get("limit"), 100), 1, 1000)

	documents, err := getAllDocuments(db, collectionName, skip, limit)
	if err != nil {
		log.Printf("Error retrieving documents: %v", err)
		http.Error(w, `{"error": {"message": "Failed to retrieve documents"}}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(documents)
}

func replaceDocumentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collectionName := r.PathValue("collection")
	if !isCollectionExists(collectionName) {
		http.Error(w, `{"error": {"message": "Collection not found"}}`, http.StatusNotFound)
		return
	}

	authToken := getAuthTokenFromRequest(r)
	if !isAuthTokenValid(authToken, collectionName, ActionReplace) {
		http.Error(w, `{"error": {"message": "Unauthorized access"}}`, http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := StoiStrict(idStr)
	if err != nil {
		http.Error(w, `{"error": {"message": "Invalid ID"}}`, http.StatusBadRequest)
		return
	}

	var document map[string]any
	err = json.NewDecoder(r.Body).Decode(&document)
	if err != nil {
		http.Error(w, `{"error": {"message": "Invalid JSON"}}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = updateDocument(db, collectionName, id, document)
	if err != nil {
		log.Printf("Error updating document: %v", err)
		http.Error(w, `{"error": {"message": "Failed to update document"}}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"message": "Document updated"}`)
}

func deleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collectionName := r.PathValue("collection")
	if !isCollectionExists(collectionName) {
		http.Error(w, `{"error": {"message": "Collection not found"}}`, http.StatusNotFound)
		return
	}

	authToken := getAuthTokenFromRequest(r)
	if !isAuthTokenValid(authToken, collectionName, ActionReplace) {
		http.Error(w, `{"error": {"message": "Unauthorized access"}}`, http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")
	id, err := StoiStrict(idStr)
	if err != nil {
		http.Error(w, `{"error": {"message": "Invalid ID"}}`, http.StatusBadRequest)
		return
	}

	err = deleteDocument(db, collectionName, id)
	if err != nil {
		log.Printf("Error deleting document: %v", err)
		http.Error(w, `{"error": {"message": "Failed to delete document"}}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"message": "Document deleted"}`)
}
