package config

import (
	"log"
	"os"
	"reflect"

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
	API_ACCESS_TOKEN string `mapstructure:"API_ACCESS_TOKEN"`

	// MaxGoRoutines is the maximum number of go routines for automated tasks / workers
	MaxGoRoutines int `mapstructure:"max_go_routines"`

	// DB :env
	DATABASE_URL string `mapstructure:"DATABASE_URL"`

	// Redis :env
	REDIS_URL string `mapstructure:"REDIS_URL"`

	// JWT :env
	JWT_SECRET string `mapstructure:"JWT_SECRET"`
}

func MustLoadConfig() *Config {
	var c = &Config{}

	v := viper.New()

	// Assuming "environment" config controls which configuration file is loaded.
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	v.SetConfigName(env)                 // Name of config file (without extension).
	v.AddConfigPath("./service/config/") // Path where config file is located.
	v.SetConfigType("yaml")              // Type of config file.

	v.AutomaticEnv() // Automatically override config values with environment variables.
	v = SetViperDefaults(v, c)

	err := v.ReadInConfig() // Read the config file.
	if err != nil {
		log.Panicf("Error reading config file, %s", err)
	}

	err = v.Unmarshal(c) // Unmarshal config into struct.
	if err != nil {
		log.Panicf("Error unmarshalling config, %s", err)
	}

	return c
}

// Due to viper's poor support for environment variables, we need to set defaults or the overrides from the actual environment variables will not work.
func SetViperDefaults(v *viper.Viper, configStruct interface{}) *viper.Viper {
	val := reflect.ValueOf(configStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure we're dealing with a struct
	if val.Kind() != reflect.Struct {
		panic("configStruct must be a struct or a pointer to a struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Use the tag `mapstructure` as the key if it exists, otherwise use the field name
		key := fieldType.Tag.Get("mapstructure")
		if key == "" {
			key = fieldType.Name
		}

		// Set the default value in Viper using the field name and its zero value
		// You might want to adjust this logic based on how you want to handle different field types
		// For simplicity, this example uses the zero value of the field's type as the default value
		v.SetDefault(key, field.Interface())
	}

	return v
}
