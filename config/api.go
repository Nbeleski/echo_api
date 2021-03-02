package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Type string `json:"type"`
		File string `json:"file"`
	} `json:"database"`

	Host string `json:"host"`
	Port string `json:"port"`
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
