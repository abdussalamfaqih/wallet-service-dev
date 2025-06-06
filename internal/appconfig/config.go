package appconfig

import (
	"log"

	"github.com/abdussalamfaqih/wallet-service-dev/pkg/config"
)

type (
	Config struct {
		App      AppConfig `yaml:"app" json:"app"`
		Database Database  `yaml:"database" json:"database"`
	}

	AppConfig struct {
		Name       string `yaml:"name" json:"name"`
		Port       string `yaml:"port" json:"port"`
		Secret_Key string `yaml:"secret_key" json:"secret_key"`
	}

	Database struct {
		Name     string `yaml:"name" json:"name"`
		Username string `yaml:"user" json:"username"`
		Password string `yaml:"pass" json:"password"`
		Host     string `yaml:"host" json:"host"`
		Port     int    `yaml:"port" json:"port"`
	}
)

func LoadConfig(path string) Config {

	var appConfig Config

	log.Println("Loading Server Configurations...")
	config.LoadConfig(&appConfig, path)
	return appConfig
}
