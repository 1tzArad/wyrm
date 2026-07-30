package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1tzArad/wyrm/internal/storage/postgres"
)

type Config struct {
	Driver string
	URL    string
}

var (
	ErrUnsupportedDriver = errors.New("Unsupported database driver")
)

func Open(cfg Config) (*sql.DB, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Open(cfg.URL)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedDriver, cfg.Driver)
	}
}

func CreateConfig(driver, url string) *Config {
	return &Config{
		Driver: driver,
		URL:    url,
	}
}
