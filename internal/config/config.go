package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type State struct {
	Config *Config
}

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(cfg)
}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not decode from file: %v", err)
	}

	return cfg, nil
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get Home directory: %v", err)
	}
	configFullPath := filepath.Join(homeDir, configFileName)
	return configFullPath, nil
}

func write(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("could not create named file %s: %v", configPath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(cfg); err != nil {
		return fmt.Errorf("could not encode config to file: %v", err)
	}
	return nil
}
