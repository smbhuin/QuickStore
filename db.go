package main

import (
	"encoding/json"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var db *sqlx.DB

type DataTable struct {
	ID        int    `db:"id"`
	CreatedAt string `db:"created_at"`
	Data      string `db:"data"`
}

func connectToDatabase(filePath string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	db, err = sqlx.Open("sqlite", filePath)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

func createSQLDDLForCollection(collectionName string) string {
	return `
	CREATE TABLE IF NOT EXISTS ` + collectionName + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		data BLOB NOT NULL
	);`
}

func migrateDatabase(db *sqlx.DB, collections []Collection) error {
	for _, collection := range collections {
		schema := createSQLDDLForCollection(collection.Name)
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}

// Store data as JSONB
func insertDocument(db *sqlx.DB, collectionName string, document map[string]any) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		return err
	}

	query := `INSERT INTO ` + collectionName + ` (data) VALUES (jsonb($1))`
	_, err = db.Exec(query, jsonData)
	return err
}

// Retrieve JSONB data
func getDocument(db *sqlx.DB, collectionName string, id int) (map[string]any, error) {
	record := DataTable{}
	query := `SELECT id, created_at, json(data) AS data FROM ` + collectionName + ` WHERE id = $1`
	err := db.Get(&record, query, id)
	if err != nil {
		return nil, err
	}

	var document map[string]any
	err = json.Unmarshal([]byte(record.Data), &document)
	if err == nil {
		document["_id"] = record.ID
		document["_created_at"] = record.CreatedAt
	}

	return document, err
}

func getAllDocuments(db *sqlx.DB, collectionName string, skip int, limit int) ([]map[string]any, error) {
	records := []DataTable{}
	query := `SELECT id, created_at, json(data) AS data FROM ` + collectionName + ` LIMIT $1 OFFSET $2`
	err := db.Select(&records, query, limit, skip)
	if err != nil {
		return nil, err
	}
	var documents []map[string]any
	for _, record := range records {
		var document map[string]any
		err = json.Unmarshal([]byte(record.Data), &document)
		if err == nil {
			document["_id"] = record.ID
			document["_created_at"] = record.CreatedAt
			documents = append(documents, document)
		}
	}
	return documents, nil

}

func deleteDocument(db *sqlx.DB, collectionName string, id int) error {
	query := `DELETE FROM ` + collectionName + ` WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func updateDocument(db *sqlx.DB, collectionName string, id int, document map[string]any) error {
	jsonData, err := json.Marshal(document)
	if err != nil {
		return err
	}

	query := `UPDATE ` + collectionName + ` SET data = jsonb($1) WHERE id = $2`
	_, err = db.Exec(query, jsonData, id)
	return err
}

func patchDocument(db *sqlx.DB, collectionName string, id int, updates map[string]any) error {
	// Retrieve existing document
	document, err := getDocument(db, collectionName, id)
	if err != nil {
		return err
	}
	for key, value := range updates {
		document[key] = value
	}
	return updateDocument(db, collectionName, id, document)
}
