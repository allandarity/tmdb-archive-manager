package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Database struct {
		Host     string `mapstructure:"host"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Port     string `mapstructure:"port"`
		Schema   string `mapstructure:"schema"`
	} `mapstructure:"database"`
	TmdbKey  string `mapstructure:"tmdb_key"`
	RestPort string `mapstructure:"rest_port"`
}

var ApplicationConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	err = viper.Unmarshal(&ApplicationConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}
}
