package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFile = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return writeConfig(*c)
}

func ReadConfig() (Config, error) {

	cfgPath, err := configPath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil

}

func configPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cfgPath := filepath.Join(homeDir, configFile)

	return cfgPath, nil
}

func writeConfig(cfg Config) error {
	cfgPath, err := configPath()
	if err != nil {
		return err
	}

	file, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}
