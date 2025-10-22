package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Environment     string `mapstructure:"ENVIRONMENT"`
	ServiceVersion  string `mapstructure:"SERVICE_VERSION"`
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	ServerPort      int    `mapstructure:"SERVER_PORT"`
	GRPCPort        int    `mapstructure:"GRPC_PORT"`
	CORSOrigins     string `mapstructure:"CORS_ORIGINS"`
	MigrationPath   string `mapstructure:"MIGRATION_PATH"`
}

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from .env file if it exists
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	v.AddConfigPath("../..")

	// Read config file (ignore error if file doesn't exist)
	_ = v.ReadInConfig()

	// Environment variables take precedence
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("ENVIRONMENT", "development")
	v.SetDefault("SERVICE_VERSION", "1.0.0")
	v.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/applicants?sslmode=disable")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("GRPC_PORT", 9090)
	v.SetDefault("CORS_ORIGINS", "*")
	v.SetDefault("MIGRATION_PATH", "internal/db/migrations")
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535")
	}

	if c.GRPCPort <= 0 || c.GRPCPort > 65535 {
		return fmt.Errorf("GRPC_PORT must be between 1 and 65535")
	}

	return nil
}

// GetCORSOrigins returns CORS origins as a slice
func (c *Config) GetCORSOrigins() []string {
	if c.CORSOrigins == "*" {
		return []string{"*"}
	}
	return strings.Split(c.CORSOrigins, ",")
}
