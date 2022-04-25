package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CFUser     string `envconfig:"cf_user"`
	CFPassword string `envconfig:"cf_password"`
	CFClient   string `envconfig:"cf_client"`
	CFSecret   string `envconfig:"cf_secret"`
}

func configLoad() (Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return Config{}, err
	}

	if config.CFUser == "" && config.CFClient == "" {
		log.Fatal("Please set CF_USER/CF_PASSWORD or CF_CLIENT/CF_SECRET")
	}

	if config.CFUser != "" && config.CFClient != "" {
		log.Println("Both CF_USER and CF_CLIENT are set. I'll use CF_CLIENT and ignore CF_USER.")
	}

	return config, nil
}
