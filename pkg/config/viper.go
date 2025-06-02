package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig(cfg interface{}) {

	viper.AddConfigPath("config")

	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
