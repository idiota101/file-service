package main

import (
	"github.com/joho/godotenv"

	"github.com/sajanjswl/file-service/pkg/cmd"

	log "github.com/sirupsen/logrus"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("No .env file found")
	}

}
func main() {
	if err := cmd.RunServer(); err != nil {
		log.Fatal(err)

	}

}
