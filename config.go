package main

import (
	"time"

	"github.com/spf13/viper"
)

const (
	// Name of the application
	AppName = "linkfire-challenge"
)

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()
	config.AutomaticEnv()

	return config
}

func (c *Config) setDefaults() {
	// Set default App configuration
	c.SetDefault("APP_ADDR", ":8080")

	// Server configuration
	c.SetDefault("SERVER_READ_TIMEOUT", 30*time.Second)
	c.SetDefault("SERVER_READ_HEADER_TIMEOUT", 10*time.Second)
	c.SetDefault("SERVER_WRITE_TIMEOUT", 30*time.Second)
}
