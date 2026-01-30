package main

import (
	"encoding/json"
	"fmt"

	"embed"
	"io/fs"
	"net/http"
)

//go:embed swagger-ui
var swagfs embed.FS

var openapiSpec []byte

func buildOpenapiSpec(config Config) ([]byte, error) {
	spec := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "QuickStore API",
			"description": "A simple document store API with authentication and schema validation",
			"version":     "1.0.0",
		},
		"servers": []map[string]any{
			{
				"url": config.OpenapiHost + "/api",
			},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
		"security": []map[string]any{
			{
				"bearerAuth": []string{},
			},
		},
	}
	specPaths := map[string]any{}
	schemas := map[string]any{}
	tags := []map[string]any{}

	schemas["ErrorResponse"] = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"error": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"message": map[string]any{
						"type": "string",
					},
				},
			},
		},
	}

	schemas["SuccessResponse"] = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type": "string",
			},
		},
	}

	specPaths["/health"] = map[string]any{
		"get": map[string]any{
			"summary": "Health check endpoint",
			"responses": map[string]any{
				"200": map[string]any{
					"description": "OK",
				},
			},
			"security": []map[string]any{},
		},
	}

	for _, collection := range config.Collections {
		schemaName := collection.Schema["title"].(string)
		schemas[schemaName] = collection.Schema
		tags = append(tags, map[string]any{
			"name":        collection.Name,
			"description": fmt.Sprintf("Operations related to the %s collection", collection.Name),
		})

		specPaths["/"+collection.Name] = map[string]any{
			"get": map[string]any{
				"summary":     "Get all documents from a collection",
				"description": "Retrieve all documents from the specified collection",
				"tags":        []string{collection.Name},
				"parameters": []map[string]any{
					{
						"name":        "skip",
						"in":          "query",
						"description": "Skip number of documents",
						"required":    false,
						"schema": map[string]any{
							"type": "integer",
						},
					},
					{
						"name":        "limit",
						"in":          "query",
						"description": "Limit number of documents",
						"required":    false,
						"schema": map[string]any{
							"type": "integer",
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Documents retrieved successfully",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"type": "array",
									"items": map[string]any{
										"$ref": fmt.Sprintf("#/components/schemas/%s", schemaName),
									},
								},
							},
						},
					},
					"401": map[string]any{
						"description": "Unauthorized access",
					},
					"404": map[string]any{
						"description": "Collection not found",
					},
				},
			},
			"post": map[string]any{
				"summary":     "Insert a new document",
				"description": "Insert a new document into the specified collection",
				"tags":        []string{collection.Name},
				"requestBody": map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"$ref": fmt.Sprintf("#/components/schemas/%s", schemaName),
							},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Document inserted successfully",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/SuccessResponse",
								},
							},
						},
					},
					"400": map[string]any{
						"description": "Invalid JSON or validation failed",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
					"401": map[string]any{
						"description": "Unauthorized access",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
					"404": map[string]any{
						"description": "Collection not found",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
		}

		specPaths[fmt.Sprintf("/%s/{id}", collection.Name)] = map[string]any{
			"get": map[string]any{
				"summary":     "Get a document by ID",
				"description": "Retrieve a specific document from the collection by ID",
				"tags":        []string{collection.Name},
				"parameters": []map[string]any{
					{
						"name":        "id",
						"in":          "path",
						"description": "Document ID",
						"required":    true,
						"schema": map[string]any{
							"type": "integer",
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Document retrieved successfully",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": fmt.Sprintf("#/components/schemas/%s", schemaName),
								},
							},
						},
					},
					"401": map[string]any{
						"description": "Unauthorized access",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
					"404": map[string]any{
						"description": "Document or collection not found",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
			"put": map[string]any{
				"summary":     "Replace a document",
				"description": "Replace an existing document with new data",
				"tags":        []string{collection.Name},
				"parameters": []map[string]any{
					{
						"name":        "id",
						"in":          "path",
						"description": "Document ID",
						"required":    true,
						"schema": map[string]any{
							"type": "integer",
						},
					},
				},
				"requestBody": map[string]any{
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{
								"$ref": fmt.Sprintf("#/components/schemas/%s", schemaName),
							},
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Document replaced successfully",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/SuccessResponse",
								},
							},
						},
					},
					"400": map[string]any{
						"description": "Invalid JSON or validation failed",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
					"401": map[string]any{
						"description": "Unauthorized access",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
					"404": map[string]any{
						"description": "Document or collection not found",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
			"delete": map[string]any{
				"summary":     "Delete a document",
				"description": "Delete a document from the collection",
				"tags":        []string{collection.Name},
				"parameters": []map[string]any{
					{
						"name":        "id",
						"in":          "path",
						"description": "Document ID",
						"required":    true,
						"schema": map[string]any{
							"type": "integer",
						},
					},
				},
				"responses": map[string]any{
					"200": map[string]any{
						"description": "Document deleted successfully",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/SuccessResponse",
								},
							},
						},
					},
					"401": map[string]any{
						"description": "Unauthorized access",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
					"404": map[string]any{
						"description": "Document or collection not found",
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{
									"$ref": "#/components/schemas/ErrorResponse",
								},
							},
						},
					},
				},
			},
		}
	}

	spec["components"].(map[string]any)["schemas"] = schemas
	spec["tags"] = tags
	spec["paths"] = specPaths

	return json.Marshal(spec)
}

// Handler returns a handler that will serve a self-hosted Swagger UI with your spec path set to /apispec.json
func SwaggerHandler() http.Handler {
	static, _ := fs.Sub(swagfs, "swagger-ui")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /apispec.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openapiSpec)
	})
	mux.Handle("/", http.FileServer(http.FS(static)))
	return mux
}
