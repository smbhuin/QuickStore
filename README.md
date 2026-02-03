# QuickStore — Lightweight Go Document Store API

A simple Go HTTP API server that allowes to store and access JSON decoments with predefined document schema.

## Running the Server

```bash
go run .
```

The server will start on `http://localhost:8080` as specified in `config.json`.

You can override the config and database paths:
```bash
go run . -config ./config.json -db ./quickstore.db
```

## Configuration

QuickStore reads `config.json` at startup and validates it against the JSON Schema in `config_schema.go`.

Example:
```json
{
  "host": "localhost",
  "port": 8080,
  "access_tokens": [
    { "name": "public", "token": "public_access_token" }
    { "name": "private", "token": "private_access_token" }
  ],
  "collections": [
    {
      "name": "products",
      "auth": {
        "all": ["private"],
        "create": [],
        "read": ["public"],
        "list": ["public"],
        "replace": [],
        "patch": [],
        "delete": []
      }
    }
  ]
}
```

Field usage:
- `host`: Hostname or IP address the server binds to.
- `port`: TCP port the server listens on.
- `access_tokens`: List of access tokens that can authenticate requests.
- `access_tokens[].name`: Friendly label used by collection auth rules.
- `access_tokens[].token`: Secret bearer token value.
- `collections`: List of collection definitions.
- `collections[].name`: Collection name used in API routes.
- `collections[].auth`: Per-action access control lists.
- `collections[].auth.all`: Tokens allowed to perform any action.
- `collections[].auth.create`: Tokens allowed to create records.
- `collections[].auth.read`: Tokens allowed to read a single record.
- `collections[].auth.list`: Tokens allowed to list records.
- `collections[].auth.replace`: Tokens allowed to replace a record.
- `collections[].auth.patch`: Tokens allowed to partially update a record.
- `collections[].auth.delete`: Tokens allowed to delete a record.
- `collections[].schema`: JSON Schema of the collection document.

## API Endpoints

- `GET /api/health` - Health check endpoint
- `GET /api/{collection}` - Get all documents from a collection
- `POST /api/{collection}` - Insert a new document into a collection
- `GET /api/{collection}/{id}` - Get a document by ID
- `PUT /api/{collection}/{id}` - Replace a document
- `DELETE /api/{collection}/{id}` - Delete a document

## API Docs

`http://localhost:8080/docs/` - Swagger UI
`http://localhost:8080/docs/apispec.json` - OpenAPI Spec JSON

## Testing the Endpoint

```bash
curl http://localhost:8080/api/health
```

Expected response:
```json
{"message": "OK"}
```

## Building

*Build in debug mode*

```bash
go build -o quickstore
./quickstore
```

*Build in release mode*

```bash
go build -o quickstore -ldflags="-s -w"
```

## Project Structure

- `go.mod` - Go module definition
- `main.go` - Main server file with the server setup and routing
- `config.go` - Configuration handling
- `db.go` - Database operations
- `routes.go` - HTTP route handlers
- `oapi.go` - OpenAPI/Swagger specification
- `README.md` - This file
