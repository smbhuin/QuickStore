package main

import (
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	return db
}

func TestCreateSQLDDLForCollection(t *testing.T) {
	collectionName := "test_collection"
	ddl := createSQLDDLForCollection(collectionName)

	expectedParts := []string{
		"CREATE TABLE IF NOT EXISTS " + collectionName,
		"id INTEGER PRIMARY KEY AUTOINCREMENT",
		"created_at DATETIME DEFAULT CURRENT_TIMESTAMP",
		"data BLOB NOT NULL",
	}

	for _, part := range expectedParts {
		if !strings.Contains(ddl, part) {
			t.Errorf("DDL does not contain expected part: %s", part)
		}
	}
}

func TestMigrateDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	collections := []Collection{
		{Name: "users"},
		{Name: "posts"},
	}

	err := migrateDatabase(db, collections)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	// Check if tables were created
	for _, collection := range collections {
		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", collection.Name)
		if err != nil {
			t.Fatalf("Failed to check table existence: %v", err)
		}
		if count != 1 {
			t.Errorf("Table %s was not created", collection.Name)
		}
	}
}

func TestInsertDocument(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	collectionName := "test_collection"
	collections := []Collection{{Name: collectionName}}
	err := migrateDatabase(db, collections)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	document := map[string]any{
		"name": "John Doe",
		"age":  30,
	}

	err = insertDocument(db, collectionName, document)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Check if document was inserted
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM "+collectionName)
	if err != nil {
		t.Fatalf("Failed to count documents: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 document, got %d", count)
	}
}

func TestGetDocument(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	collectionName := "test_collection"
	collections := []Collection{{Name: collectionName}}
	err := migrateDatabase(db, collections)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	document := map[string]any{
		"name": "Jane Doe",
		"age":  25,
	}

	err = insertDocument(db, collectionName, document)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Get the document
	retrieved, err := getDocument(db, collectionName, 1)
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}

	if retrieved["name"] != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got %v", retrieved["name"])
	}
	if retrieved["age"] != float64(25) { // JSON unmarshals numbers as float64
		t.Errorf("Expected age 25, got %v", retrieved["age"])
	}
	if retrieved["_id"] != 1 {
		t.Errorf("Expected _id 1, got %v", retrieved["_id"])
	}
	if retrieved["_created_at"] == nil {
		t.Error("Expected _created_at to be set")
	}
}

func TestGetDocumentNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	collectionName := "test_collection"
	collections := []Collection{{Name: collectionName}}
	err := migrateDatabase(db, collections)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	_, err = getDocument(db, collectionName, 999)
	if err == nil {
		t.Error("Expected error for non-existent document")
	}
}
