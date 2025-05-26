package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return fmt.Errorf("Error loading .env file:", err)
	}

	return nil
}

type EvmConfig interface {
}
