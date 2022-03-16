package main

import (
	"github.com/joho/godotenv"
	"github.com/meghashyamc/country-incomes/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	godotenv.Load()
	cmd.Execute()
}
