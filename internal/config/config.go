package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig
	App    AppConfig
}

type ServerConfig struct {
	Port string
}

type AppConfig struct {
	Environment string
	MaxLimit    int
}

func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("MAX_LIMIT", 10000)

	// Bind environment variables
	viper.AutomaticEnv()

	// Build config
	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
		},
		App: AppConfig{
			Environment: viper.GetString("ENVIRONMENT"),
			MaxLimit:    viper.GetInt("MAX_LIMIT"),
		},
	}

	return config, nil
}