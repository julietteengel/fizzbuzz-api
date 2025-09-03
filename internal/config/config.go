package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	App      AppConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type AppConfig struct {
	Environment string
	MaxLimit    int
}

type DatabaseConfig struct {
	URL          string
	StatsStorage string
}

func Load() (*Config, error) {
	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("MAX_LIMIT", 10000)
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/fizzbuzz_db?sslmode=disable")
	viper.SetDefault("STATS_STORAGE", "memory")

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
		Database: DatabaseConfig{
			URL:          viper.GetString("DATABASE_URL"),
			StatsStorage: viper.GetString("STATS_STORAGE"),
		},
	}

	return config, nil
}