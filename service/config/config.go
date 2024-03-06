package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	// API
	ServiceName string `mapstructure:"service_name"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`

	// Logging
	LogLevel string `mapstructure:"log_level"`
	LogFile  string `mapstructure:"log_file"`

	// Client :env
	APIAccessToken string `mapstructure:"st_api_access_token"`

	// MaxGoRoutines is the maximum number of go routines for automated tasks / workers
	MaxGoRoutines int `mapstructure:"max_go_routines"`

	// DB :env
	DatabaseURL string `mapstructure:"database_url"`

	// Redis :env
	RedisURL string `mapstructure:"redis_url"`

	// JWT :env
	JWTSecret string `mapstructure:"jwt_secret"`
}

func MustLoadConfig() *Config {
	var c *Config
	var env string

	v := viper.New()

	if env = os.Getenv("ENVIRONMENT"); env == "" {
		env = "development"
	}

	v.SetConfigName(env)
	v.AddConfigPath("./service/config/")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		log.Panicf("Error reading config file, %s", err.Error())
	}

	err = v.Unmarshal(&c)
	if err != nil {
		log.Panicf("Error unmarshalling config, %s", err.Error())
	}

	return c
}
