// Package env loads and parses environment variables from the .env file on the root directory of the app
package env

import (
	"log"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	Port       uint   `env:"PORT"`
	Connect    string `env:"CONNECT"`
	Disconnect string `env:"DISCONNECT"`
	Exit       string `env:"EXIT"`
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
