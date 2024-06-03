package config

import (
	"github.com/joeshaw/envdecode"
	"log"
)

type Conf struct {
	Server ConfServer
}

type ConfServer struct {
	Port  int  `env:"SERVER_PORT"`
	Debug bool `env:"SERVER_DEBUG"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
