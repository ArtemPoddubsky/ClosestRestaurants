package main

import (
	"main/internal/app"
	"main/internal/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("Building")

	cfg := config.GetConfig()
	application := app.NewApp(&cfg)

	application.Run()
}
