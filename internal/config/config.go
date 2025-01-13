package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"errors"
)

const configFileName = ".gatorconfig.json"


type Config struct {
	DBURL 			string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to retrieve users home directory: " + err.Error())
	}
	return filepath.Join(home, configFileName), nil 
}


func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, errors.New("failed to retrieve config file path: " + err.Error())
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, errors.New("failed to read config file: " + err.Error())
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, errors.New("failed to unmarshal json from config file: " + err.Error())
	}
	return cfg, nil
}


func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return errors.New("failed to retrieve config file path: " + err.Error())
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return errors.New("failed to marshal json: " + err.Error())
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return errors.New("failed to write config file: " + err.Error())
	}
	return nil
}
