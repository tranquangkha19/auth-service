package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppName  string         `mapstructure:"APP_NAME"`
	Port     string         `mapstructure:"PORT"`
	Database DatabaseConfig `mapstructure:"DATABASE"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"SSL_MODE"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	// Set environment variable mappings for database credentials
	viper.BindEnv("DATABASE.USER", "DB_USER")
	viper.BindEnv("DATABASE.PASSWORD", "DB_PASSWORD")
	viper.BindEnv("DATABASE.HOST", "DB_HOST")
	viper.BindEnv("DATABASE.PORT", "DB_PORT")
	viper.BindEnv("DATABASE.DB_NAME", "DB_NAME")
	viper.BindEnv("DATABASE.SSL_MODE", "SSL_MODE")

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
