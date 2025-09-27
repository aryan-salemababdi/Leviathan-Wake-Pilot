package config

import (
	"fmt"
	"os"
)

type Config struct {
	KeyDBAddress      string
	GrpcServerPort    string
	ExchangeApiKey    string
	ExchangeApiSecret string
}

func Load() *Config {
	keydbHost := os.Getenv("KEYDB_HOST")
	if keydbHost == "" {
		keydbHost = "localhost"
	}

	keydbPort := os.Getenv("KEYDB_PORT")
	if keydbPort == "" {
		keydbPort = "6379"
	}

	grpcPort := os.Getenv("GRPC_SERVER_PORT")
	if grpcPort == "" {
		grpcPort = ":50051"
	}

	return &Config{
		KeyDBAddress:      fmt.Sprintf("%s:%s", keydbHost, keydbPort),
		GrpcServerPort:    grpcPort,
		ExchangeApiKey:    os.Getenv("EXCHANGE_API_KEY"),
		ExchangeApiSecret: os.Getenv("EXCHANGE_API_SECRET"),
	}
}
