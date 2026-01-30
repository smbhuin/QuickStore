package main

import (
	"encoding/json"
	"os"

	"github.com/go-playground/validator/v10"
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
	Host         string        `json:"host" validate:"required"`
	OpenapiHost  string        `json:"openapi_host" validate:"required"`
	Port         int           `json:"port" validate:"required"`
	AccessTokens []AccessToken `json:"access_tokens" validate:"required"`
	Collections  []Collection  `json:"collections" validate:"required"`
}

type AccessToken struct {
	Name  string `json:"name" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type Collection struct {
	Name   string         `json:"name" validate:"required"`
	Auth   CollectionAuth `json:"auth" validate:"required"`
	Schema map[string]any `json:"schema" validate:"required"`
}

type CollectionAuth struct {
	All     []string `json:"all" validate:"required"`
	Create  []string `json:"create" validate:"required"`
	Read    []string `json:"read" validate:"required"`
	List    []string `json:"list" validate:"required"`
	Replace []string `json:"replace" validate:"required"`
	Patch   []string `json:"patch" validate:"required"`
	Delete  []string `json:"delete" validate:"required"`
}

func readConfig(fileName string) (Config, error) {
	var config Config
	file, err := os.Open(fileName)
	if err != nil {
		return config, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	// Add validation
	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getCollectionByName(collectionName string) *Collection {
	for _, collection := range config.Collections {
		if collection.Name == collectionName {
			return &collection
		}
	}
	return nil
}
