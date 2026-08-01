package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/1tzArad/wyrm/internal/config"
	"github.com/1tzArad/wyrm/internal/storage/postgres"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func init() {
	newLog := log.NewWithOptions(os.Stdout, log.Options{
		TimeFormat:      time.DateTime,
		Prefix:          "WYRM",
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	log.SetDefault(newLog)
}

func main() {
	cfg := config.Setup()

	src, err := iofs.New(postgres.MigrationsFS, "migrations")
	if err != nil {
		log.Fatal("Failed to load migration files!", "err", err.Error())
		return
	}

	dbURL := url.URL{
		Scheme: "pgx5",
		User:   url.UserPassword(cfg.PostgresDB.User, cfg.PostgresDB.Pass),
		Host:   fmt.Sprintf("%s:%d", cfg.PostgresDB.Host, cfg.PostgresDB.Port),
		Path:   "/" + cfg.PostgresDB.Name,
	}
	query := dbURL.Query()
	query.Set("sslmode", cfg.PostgresDB.SSLMode)
	dbURL.RawQuery = query.Encode()

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL.String())
	if err != nil {
		log.Fatal("Failed to initialize migrations!", "err", err.Error())
		return
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No pending migrations")
			return
		}
		log.Fatal("Failed to run migrations!", "err", err.Error())
		return
	}

	log.Info("Database migrations completed successfully")
}
