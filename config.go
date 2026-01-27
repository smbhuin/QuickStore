package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host         string        `json:"host"`
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
	Name   string         `json:"name"`
	Schema map[string]any `json:"schema"`
}

func readConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
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
