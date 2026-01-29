package main

import (
	"encoding/json"
	"os"
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
	Debug        bool          `json:"debug"`
	AccessTokens []AccessToken `json:"access_tokens"`
	Collections  []Collection  `json:"collections"`
}

type AccessToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type Collection struct {
	Name   string              `json:"name"`
	Auth   map[string][]string `json:"auth"`
	Schema map[string]any      `json:"schema"`
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
