package main

import (
	"os"
)

type config struct {
	apiKey string
}

func getConfig() {
	// Set defaults

	// Get the config from file

	// Then from env
	readEnvConfig()

	// Finally from command options
}

func readEnvConfig() {
	key := os.Getenv("TODOISTA_KEY")

	if key != "" {
		app.cfg.apiKey = key
	}
}
