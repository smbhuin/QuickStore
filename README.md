# QuickStore - Go HTTP Server

A simple Go HTTP server with a document store API.

## Running the Server

```bash
go run .
```

The server will start on `http://localhost:8080`

## API Endpoints

- `GET /api/health` - Health check endpoint
- `GET /api/{collection}` - Get all documents from a collection
- `POST /api/{collection}` - Insert a new document into a collection
- `GET /api/{collection}/{id}` - Get a document by ID
- `PUT /api/{collection}/{id}` - Replace a document
- `DELETE /api/{collection}/{id}` - Delete a document

## API Docs

`http://localhost:8080/apispec/` - Swagger UI
`http://localhost:8080/apispec/apispec.json` - OpenAPI Spec JSON

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
