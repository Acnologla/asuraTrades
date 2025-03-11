package config

type (
	PostgresConfig struct {
		DatabaseURL string
	}

	JWTConfig struct {
		Secret string
	}

	Config struct {
		Postgres PostgresConfig
		JWT      JWTConfig
	}
)
