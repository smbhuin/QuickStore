package main

import (
	"github.com/xeipuuv/gojsonschema"
)

var schemaCache map[string]gojsonschema.JSONLoader

func buildSchemaCache(collections []Collection) map[string]gojsonschema.JSONLoader {
	var schemaCache = make(map[string]gojsonschema.JSONLoader)
	for _, collection := range collections {
		schemaCache[collection.Name] = gojsonschema.NewGoLoader(collection.Schema)
	}
	return schemaCache
}

func validateJSON(json map[string]any, schemaLoader gojsonschema.JSONLoader) bool {

	// Create document loader from the json data
	documentLoader := gojsonschema.NewGoLoader(json)

	// Validate the json against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false
	}

	return result.Valid()
}

func validateJSONByCollectionName(json map[string]any, collectionName string) bool {
	schemaLoader, exists := schemaCache[collectionName]

	if !exists {
		return false
	}

	return validateJSON(json, schemaLoader)
}

func isCollectionExists(collectionName string) bool {
	_, exists := schemaCache[collectionName]
	return exists
}
