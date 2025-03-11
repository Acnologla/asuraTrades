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

	HTTPConfig struct {
		Port                  string
		GenerateTokenPassword string
	}

	Config struct {
		PostgresConfig PostgresConfig
		JWTConfig      JWTConfig
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
	httpConfig := HTTPConfig{
		GenerateTokenPassword: os.Getenv("GENERATE_TOKEN_PASSWORD"),
		Port:                  os.Getenv("PORT"),
	}

	return &Config{
		PostgresConfig: postgresConfig,
		HTTPConfig:     httpConfig,
		JWTConfig:      jwtConfig,
		Production:     isProduction,
	}, nil
}
