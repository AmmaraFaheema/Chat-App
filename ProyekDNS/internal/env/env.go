package env

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	Port int `env:"PORT"`
}

func New() *Env {
	log.Println("loading environment variables")

	err := godotenv.Load()
	if err != nil {
		log.Panic("failed to load environment variables")
	}

	envParsed := new(Env)
	err = env.Parse(envParsed)
	if err != nil {
		log.Panic("failed to parse environment variables")
	}

	return envParsed
}
