package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ListenAddrStrings []string `json:"listenAddrStrings"`
	WhitelistPeers    []string `json:"whitelistPeers"`
	WhitelistAddrs    []string `json:"whitelistAddrs"`
}

func (config Config) String() string {
	b, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func loadConfig() (Config, error) {
	var config Config
	data, err := os.ReadFile("config.json")
	if err != nil {
		return config, fmt.Errorf("Error loading Config: %w", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("Error parsing Config: %w", err)
	}
	return config, nil
}
