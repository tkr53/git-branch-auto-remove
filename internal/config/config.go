package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the application's configuration.
type Config struct {
	ProtectedBranches []string `mapstructure:"protected_branches"`
}

// ConfigLoader defines the interface for loading configuration.
type ConfigLoader interface {
	LoadConfig() (*Config, error)
}

// ViperConfigLoader implements ConfigLoader using Viper.
type ViperConfigLoader struct{}

// LoadConfig reads configuration from .gbar.yml.
func (l *ViperConfigLoader) LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Set default values
	viper.SetDefault("protected_branches", []string{"main", "master", "develop"})

	// Set config file name and type
	viper.SetConfigName(".gbar")
	viper.SetConfigType("yaml")

	// Add config paths: current directory and user's home directory
	viper.AddConfigPath(".")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and use defaults
			fmt.Println("No config file found, using default settings.")
		} else {
			// Config file was found but another error was encountered
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}