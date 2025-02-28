package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // config file name without extension
	viper.SetConfigType("yaml")   // or viper.SetConfigType("YAML")

	viper.AddConfigPath(".")        // looking for config in the working directory
	viper.AddConfigPath("./config") // looking for config in ./config/

	// Set default values
	setDefaults()

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired
		fmt.Println("No config file found. Using environment variables and defaults.")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "postgres")
	viper.SetDefault("database.sslmode", "disable")

	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.readtimeout", 15)  // seconds
	viper.SetDefault("server.writetimeout", 15) // seconds
	viper.SetDefault("server.idletimeout", 60)  // seconds
}

// GetDSN returns the PostgreSQL DSN string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
