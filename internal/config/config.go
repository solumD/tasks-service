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
	loggerLevelEnv    = "LOGGER_LEVEL"
)

// Config конфиг
type Config struct {
	httpServerHost string
	httpServerPort string
	loggerLevel    string
}

// ServerAddr возвращает адрес сервера
func (cfg *Config) ServerAddr() string {
	return net.JoinHostPort(cfg.httpServerHost, cfg.httpServerPort)
}

// LoggerLevel возвращает уровень логирования
func (c *Config) LoggerLevel() string {
	return c.loggerLevel
}

// MustLoad загружает конфиг из файла .env
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

	loggerLevel := os.Getenv(loggerLevelEnv)
	if len(loggerLevel) == 0 {
		log.Fatal("logger level not found")
	}

	return &Config{
		httpServerHost: serverHost,
		httpServerPort: serverPort,
		loggerLevel:    loggerLevel,
	}
}
