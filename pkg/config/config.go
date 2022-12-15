package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

// SyncatDBConfig is the configuration for the database connection
type SyncatDBConfig struct {
	// Database filename
	Filename string `yaml:"filename"`
}

// SyncatSyncConfig is the configuration for syncing
type SyncatSyncConfig struct {
	// Directories to sync
	Directories []string `yaml:"directories"`
}

type SyncatProtocolConfig struct {
	// Buffer size for TCP server
	BufferSize int `yaml:"buffer_size"`
	// Timeout configuration for TCP server
	Timeout int `yaml:"timeout"`
	// Ping interval for TCP server
	PingInterval int `yaml:"ping_interval"`
}

type SyncatAuthConfig struct {
	// Token for authentication
	Token string `yaml:"token"`
}

type SyncatConfig struct {
	// Database configuration
	Db SyncatDBConfig `yaml:"db"`
	// Sync configuration
	Sync SyncatSyncConfig `yaml:"sync"`
	// Protocol configuration
	Protocol SyncatProtocolConfig `yaml:"protocol"`
	// Authentication configuration
	Auth SyncatAuthConfig `yaml:"auth"`
}

var config SyncatConfig

func LoadConfig() error {
	// Obtain the executable path
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	exPath := filepath.Dir(ex)
	// Obtain config file path
	configPath := filepath.Join(exPath, "../config/config.yml")
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
	if err != nil {
		return err
	}
	// Unmarshal config file
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return err
	}
	// Obtain config file path
	config.Db.Filename = filepath.Join(exPath, "..", config.Db.Filename)
	for i := range config.Sync.Directories {
		config.Sync.Directories[i] = filepath.Join(exPath, "..", config.Sync.Directories[i])
	}
	return nil
}

func GetConfig() SyncatConfig {
	return config
}
