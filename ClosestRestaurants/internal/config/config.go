package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config stores all information from config file.
type Config struct {
	DB struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
	} `toml:"database"`
}

// GetConfig reads configuration file and stores it in Config.
func GetConfig() Config {
	var cfg Config

	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.SetConfigType("toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("viper.ReadInConfig: ", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalln("viper.Unmarshal:", err)
	}

	return cfg
}
