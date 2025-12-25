package config

import (
	"log"
	"net"
	"os"

	"github.com/solumD/tasks-service/pkg/env"
)

const (
	configPath = ".env"

	httpServerHostEnv = "HTTP_SERVER_HOST"
	httpServerPortEnv = "HTTP_SERVER_PORT"
)

type Config struct {
	httpServerHost string `env:"HTTP_SERVER_HOST"`
	httpServerPort string `env:"HTTP_SERVER_PORT"`
}

func (cfg *Config) ServerAddr() string {
	return net.JoinHostPort(cfg.httpServerHost, cfg.httpServerPort)
}

func MustLoad() *Config {
	err := env.LoadEnv(configPath)
	if err != nil {
		log.Fatalf("failed to load config from %s: %v", configPath, err)
	}

	serverHost := os.Getenv(httpServerHostEnv)
	if len(serverHost) == 0 {
		log.Fatal("server host not found")
	}

	serverPort := os.Getenv(httpServerPortEnv)
	if len(serverPort) == 0 {
		log.Fatal("server port not found")
	}

	return &Config{
		httpServerHost: serverHost,
		httpServerPort: serverPort,
	}
}
