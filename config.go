package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

var config Config

const (
	ActionAll     = "all"
	ActionCreate  = "create"
	ActionRead    = "read"
	ActionList    = "list"
	ActionReplace = "replace"
	ActionPatch   = "patch"
	ActionDelete  = "delete"
)

type Config struct {
	Host         string        `json:"host"`
	OpenapiHost  string        `json:"openapi_host"`
	Port         int           `json:"port"`
	AccessTokens []AccessToken `json:"access_tokens"`
	Collections  []Collection  `json:"collections"`
}

type AccessToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Collection struct {
	Name   string         `json:"name"`
	Auth   CollectionAuth `json:"auth"`
	Schema map[string]any `json:"schema"`
}

type CollectionAuth struct {
	All     []string `json:"all"`
	Create  []string `json:"create"`
	Read    []string `json:"read"`
	List    []string `json:"list"`
	Replace []string `json:"replace"`
	Patch   []string `json:"patch"`
	Delete  []string `json:"delete"`
}

func readConfig(fileName string) (Config, error) {
	var config Config
	content, err := os.ReadFile(fileName)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(content, &config)
	if err != nil {
		return config, err
	}

	documentLoader := gojsonschema.NewStringLoader(string(content))
	schemaLoader := gojsonschema.NewStringLoader(configSchema)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return config, err
	}
	if result.Valid() {
		return config, nil
	}
	log.Printf("Found error in config. see errors :\n")
	for _, err := range result.Errors() {
		// Err implements the ResultError interface
		log.Printf("- %s\n", err)
	}
	return config, errors.New("Config is invalid")
}

func getCollectionByName(collectionName string) *Collection {
	for _, collection := range config.Collections {
		if collection.Name == collectionName {
			return &collection
		}
	}
	return nil
}
