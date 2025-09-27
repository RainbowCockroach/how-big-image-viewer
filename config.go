package main

import (
	"encoding/json"
	"os"
)

type Size struct {
	Name   string `json:"name"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type Config struct {
	Scale float64 `json:"scale"`
	Sizes []Size  `json:"sizes"`
}

func loadConfig() (*Config, error) {
	config := &Config{
		Scale: 1.0,
		Sizes: []Size{},
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("config.json", data, 0644)
}