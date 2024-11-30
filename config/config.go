package config

import (
	"log"

	"github.com/spf13/viper"
)

const (
	Port    = ":9000"
	Host    = "localhost"
	Version = "v1.0.1"
)

func ReadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error in ReadConfig")
	}
}
