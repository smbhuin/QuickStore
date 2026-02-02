package main

var configSchema string = `
{
  "type": "object",
  "properties": {
    "host": {
      "type": "string"
    },
    "port": {
      "type": "integer"
    },
    "access_tokens": {
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "token": {
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
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "auth": {
              "type": "object",
              "properties": {
                "all": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "create": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "read": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "list": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "replace": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "patch": {
                  "type": "array",
                  "items": [
                    {
                      "type": "string"
                    }
                  ]
                },
                "delete": {
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
            }
          },
          "required": [
            "name",
            "auth"
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
