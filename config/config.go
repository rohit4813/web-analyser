package config

import (
	"github.com/joeshaw/envdecode"
	"log"
	"time"
)

type Conf struct {
	Server ServerConf
	Client ClientConf
}

// ServerConf is a struct for the server configurations
type ServerConf struct {
	Port         int           `env:"SERVER_PORT"`
	Debug        bool          `env:"SERVER_DEBUG"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE"`
}

// ClientConf is a struct for the client configurations
type ClientConf struct {
	Timeout time.Duration `env:"CLIENT_TIMEOUT"`
}

// New maps the environment variables to Conf using envdecode pkg
func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
