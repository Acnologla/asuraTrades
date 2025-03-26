package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	PostgresConfig struct {
		DatabaseURL string
	}

	JWTConfig struct {
		Secret string
	}

	GrpcConfig struct {
		Port  string
		Token string
	}

	HTTPConfig struct {
		Port                  string
		GenerateTokenPassword string
		ProductionURL         string
	}

	Config struct {
		PostgresConfig PostgresConfig
		JWTConfig      JWTConfig
		GrpcConfig     GrpcConfig
		HTTPConfig     HTTPConfig
		Production     bool
	}
)

func LoadConfig() (*Config, error) {
	isProduction := os.Getenv("PRODUCTION") == "TRUE"
	if !isProduction {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	jwtConfig := JWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
	}
	postgresConfig := PostgresConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
	grpcConfig := GrpcConfig{
		Port:  os.Getenv("GRPC_PORT"),
		Token: os.Getenv("GRPC_TOKEN"),
	}
	httpConfig := HTTPConfig{
		GenerateTokenPassword: os.Getenv("GENERATE_TOKEN_PASSWORD"),
		Port:                  os.Getenv("PORT"),
		ProductionURL:         os.Getenv("PRODUCTION_DOMAIN"),
	}

	return &Config{
		PostgresConfig: postgresConfig,
		HTTPConfig:     httpConfig,
		JWTConfig:      jwtConfig,
		Production:     isProduction,
		GrpcConfig:     grpcConfig,
	}, nil
}
