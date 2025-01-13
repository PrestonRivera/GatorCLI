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

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, errors.New("failed to open config file: " + err.Error())
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, errors.New("failed to parse json from config file: " + err.Error())
	}
	return cfg, nil 
}


func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return errors.New("failed to retrieve config file path: " + err.Error())
	}

	file, err := os.Create(configPath)
	if err != nil {
		return errors.New("failed to open file for  writing: " + err.Error())
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return errors.New("failed to encode json to config file: " + err.Error())
	}
	return nil
}


func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)
}