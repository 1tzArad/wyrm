package config

type Config struct {
	JWTSecret string
}

func Setup() *Config {
	return &Config{
		JWTSecret: GetEnv("JWT_SECRET"),
	}
}
