package config

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret  string
	PostgresDB *PostgresConfig
}

type PostgresConfig struct {
	User    string
	Pass    string
	Host    string
	Port    int
	Name    string
	SSLMode string
}

func Setup() *Config {
	setupDotEnv()
	postgresDBPortStr := GetEnv("POSTGRES_PORT")
	postgresDBPort, err := strconv.Atoi(postgresDBPortStr)
	log.Info(postgresDBPortStr)
	if err != nil {
		log.Fatal("POSTGRES_PORT must be a number!")
		return nil
	}

	postgresSSLMode := GetEnv("POSTGRES_SSL_MODE")

	switch postgresSSLMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		log.Fatalf("invalid POSTGRES_SSL_MODE: %q", postgresSSLMode)
	}

	return &Config{
		JWTSecret: GetEnv("JWT_SECRET"),
		PostgresDB: &PostgresConfig{
			User:    GetEnv("POSTGRES_USER"),
			Pass:    GetEnv("POSTGRES_PASS"),
			Host:    GetEnv("POSTGRES_HOST"),
			Port:    postgresDBPort,
			Name:    GetEnv("POSTGRES_NAME"),
			SSLMode: postgresSSLMode,
		},
	}
}

func (cfg *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Pass,
		cfg.Name,
		cfg.SSLMode,
	)
}

func setupDotEnv() {
	godotenv.Load()
}
