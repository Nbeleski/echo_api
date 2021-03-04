package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Driver     string `json:"driver"`
		Datasource string `json:"datasource"`
	} `json:"database"`

	Host       string `json:"host"`
	Port       string `json:"port"`
	JWT_secret string `json:"jwt_secret"`
	MapsAPIKey string `json:"maps_apikey"`
}

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, nil
}
