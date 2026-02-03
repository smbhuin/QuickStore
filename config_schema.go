package main

var configSchema string = `
{
  "type": "object",
  "properties": {
    "host": {
      "description": "Hostname or IP address the server binds to.",
      "type": "string"
    },
    "port": {
      "description": "TCP port the server listens on.",
      "type": "integer"
    },
    "access_tokens": {
      "description": "List of access tokens that can authenticate requests.",
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "description": "Human-friendly label for the access token.",
              "type": "string"
            },
            "token": {
              "description": "Secret bearer token used for authentication.",
              "type": "string"
            }
          },
          "required": [
            "name",
            "token"
          ]
        }
      ]
    },
    "collections": {
      "description": "Collection definitions that configure data storage and access control.",
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "description": "Unique collection name used in API routes.",
              "type": "string"
            },
            "auth": {
              "description": "Per-action access control lists for the collection.",
              "type": "object",
              "properties": {
                "all": {
                  "description": "Tokens allowed to perform any action on the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "create": {
                  "description": "Tokens allowed to create records in the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "read": {
                  "description": "Tokens allowed to read a single record in the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "list": {
                  "description": "Tokens allowed to list records in the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "replace": {
                  "description": "Tokens allowed to replace an entire record in the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "patch": {
                  "description": "Tokens allowed to partially update a record in the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "delete": {
                  "description": "Tokens allowed to delete a record in the collection.",
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                }
              },
              "required": [
                "all",
                "create",
                "read",
                "list",
                "replace",
                "patch",
                "delete"
              ]
            },
            "schema": {
              "type": "object"
            }
          },
          "required": [
            "name",
            "auth",
            "schema"
          ]
        }
      ]
    }
  },
  "required": [
    "host",
    "port",
    "access_tokens",
    "collections"
  ]
}
`
