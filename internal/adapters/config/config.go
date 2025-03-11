package config

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
		Postgres PostgresConfig
		JWT      JWTConfig
		HTTP     HTTPConfig
	}
)
