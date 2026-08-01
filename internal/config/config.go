package config

import (
	"strconv"

	"github.com/charmbracelet/log"
)

type Config struct {
	JWTSecret  string
	PostgresDB *PostgresConfig
}

type PostgresConfig struct {
	User string
	Pass string
	Host string
	Port int
	Name string
}

func Setup() *Config {
	postgres_db_port_str := GetEnv("POSTGRES_PORT")
	postgres_db_port, err := strconv.Atoi(postgres_db_port_str)
	if err != nil {
		log.Fatal("POSTGRES_PORT must be a number!")
		return nil
	}
	return &Config{
		JWTSecret: GetEnv("JWT_SECRET"),
		PostgresDB: &PostgresConfig{
			User: GetEnv("POSTGRES_USER"),
			Pass: GetEnv("POSTGRES_PASS"),
			Host: GetEnv("POSTGRES_HOST"),
			Port: postgres_db_port,
			Name: GetEnv("POSTGRES_NAME"),
		},
	}
}
