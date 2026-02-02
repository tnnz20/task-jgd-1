package config

import (
	"log" 

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port        string
	Environment string
	LogLevel    string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	PoolMode string
}

// NewViper creates and returns a new Viper instance with environment variables
func NewViper() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()

	v.AddConfigPath(".")
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		// Check if the error is a "Config File Not Found" error
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// If it's a different error (e.g., syntax error in .env), log it
			log.Printf("Warning: .env file found but could not be parsed: %v", err)
		}
	}
	return v
}

// NewConfig loads configuration from viper instance
func NewConfig(v *viper.Viper) *Config {
	config := &Config{
		App: AppConfig{
			Port:        v.GetString("PORT"),
			Environment: v.GetString("ENVIRONMENT"),
			LogLevel:    v.GetString("LOG_LEVEL"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetString("DB_PORT"),
			Name:     v.GetString("DB_NAME"),
			User:     v.GetString("DB_USER"),
			Password: v.GetString("DB_PASSWORD"),
			PoolMode: v.GetString("DB_POOLMODE"),
		},
	}

	return config
}
