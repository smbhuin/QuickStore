# QuickStore - Go HTTP Server

A simple Go HTTP server with a hello world endpoint.

## Running the Server

```bash
go run .
```

The server will start on `http://localhost:8080`

## Testing the Endpoint

```bash
curl http://localhost:8080
```

Expected response:
```json
{"message": "Hello World"}
```

## Building

```bash
go build -o quickstore
./quickstore
```

## Project Structure

- `go.mod` - Go module definition
- `main.go` - Main server file with the hello world endpoint
- `README.md` - This file
