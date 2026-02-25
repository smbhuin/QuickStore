# multi-stage build for the Go application
FROM golang:1.25-alpine AS builder

# ensure modules are downloaded first to leverage build cache
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# copy source and build binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o quickstore .

# final minimal image
FROM scratch

# copy binary and necessary assets
COPY --from=builder /app/quickstore /quickstore
# optional: copy sample config so users can override
COPY --from=builder /app/config.example.json /config.example.json

# default working directory inside container
WORKDIR /

# expose port defined in config (default 8080)
EXPOSE 8080

# run the binary by default; allow overriding args
ENTRYPOINT ["/quickstore"]
