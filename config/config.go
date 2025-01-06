package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	EtherscanApi = `api.etherscan.io`
)

var (
	API string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	API = os.Getenv("API")
	if API == "" {
		log.Fatal("Error loading API KEY:", err)
	}
}
