package server

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

// SyncatServerConfig is the configuration for the Syncat server
type SyncatServerConfig struct {
	// Port for server
	Port int `yaml:"port"`
	// Host for server
	Host string `yaml:"host"`
}

var serverConfig SyncatServerConfig

func LoadConfig() error {
	// Obtain the executable path
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)
	// Obtain config file path
	configPath := filepath.Join(exPath, "../config/server_config.yml")
	// Open config file
	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)
	// Read config file
	configBytes, err := io.ReadAll(configFile)
	// Unmarshal config file
	err = yaml.Unmarshal(configBytes, &serverConfig)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() SyncatServerConfig {
	return serverConfig
}
