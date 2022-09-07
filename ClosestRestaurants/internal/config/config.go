package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Db struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Database string `toml:"database"`
	} `toml:"database"`
}

func GetConfig() Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	return cfg
}
