package config

type (
	PostgresConfig struct {
		DatabaseURL string
	}

	JWTConfig struct {
		Secret string
	}

	HTTPConfig struct {
		Port string
	}

	Config struct {
		Postgres PostgresConfig
		JWT      JWTConfig
	}
)
