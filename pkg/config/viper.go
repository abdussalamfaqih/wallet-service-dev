package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig(cfg interface{}, path string) {

	viper.AddConfigPath("config")

	paths := strings.Split(path, ".")
	viper.SetConfigName(paths[0])
	viper.SetConfigType(paths[1])

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
