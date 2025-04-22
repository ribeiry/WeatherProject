package config

import (
	"log"
	"os"

	"github.com/subosito/gotenv"
)

var (
	BaseURL      string
	SynopticPath string
	APIKey       string
	Endpoint     string
	DBConnection string
)

func LoadEnv() {

	err := gotenv.Load("resources/resources.env")
	dir, _ := os.Getwd()

	log.Println("Diretorio raiz do projeto é: ", dir)

	if err != nil {
		log.Println("erro ao carregar o .env")
		log.Println(err)
	}

	BaseURL = os.Getenv("BASE_URL")
	SynopticPath = os.Getenv("SYNOPTIC")
	APIKey = os.Getenv("APIKEY")
	Endpoint = os.Getenv("ENDPOINT")
	DBConnection = os.Getenv("DB_CONNECTION")
}
