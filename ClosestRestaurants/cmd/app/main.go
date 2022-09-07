package main

import (
	log "github.com/sirupsen/logrus"
	"main/internal/app"
	"main/internal/config"
)

func main() {
	log.Infoln("Building")
	cfg := config.GetConfig()
	app.NewApp(cfg).Run()
}
