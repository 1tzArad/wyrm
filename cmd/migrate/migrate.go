package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/1tzArad/wyrm/internal/config"
	"github.com/1tzArad/wyrm/internal/storage/postgres"
	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
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

const usageText = `Usage: migrate <command> [arguments]

Commands:
  up                  Apply all pending migrations
  down                Roll back the latest migration
  steps <n>           Apply (n > 0) or roll back (n < 0) n migrations
  version             Show the current migration version and dirty state
  force <version>     Force set the migration version
  drop [-f]           Drop all database objects created by migrations
  status              Show migration status
  help                Show this help message
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage(os.Stderr)
		os.Exit(2)
	}

	cmd := args[0]
	rest := args[1:]

	switch cmd {
	case "up", "down", "version", "status":
		if len(rest) != 0 {
			usageError(cmd, "unexpected arguments")
		}
	case "steps", "force":
		if len(rest) != 1 {
			usageError(cmd, "expected exactly one argument")
		}
	case "drop":
		if len(rest) > 1 {
			usageError(cmd, "unexpected arguments")
		}
		if len(rest) == 1 && rest[0] != "-f" && rest[0] != "--force" {
			usageError(cmd, fmt.Sprintf("unknown argument %q", rest[0]))
		}
	case "help", "-h", "--help":
		usage(os.Stdout)
		return
	default:
		usageError(cmd, fmt.Sprintf("unknown command %q", cmd))
	}

	m, src, err := openMigrate()
	if err != nil {
		log.Fatal("Failed to open database!", "err", err.Error())
		return
	}
	defer m.Close()

	switch cmd {
	case "up":
		runUp(m)
	case "down":
		runDown(m)
	case "steps":
		runSteps(m, atoi(cmd, rest[0]))
	case "version":
		runVersion(m)
	case "force":
		runForce(m, atoi(cmd, rest[0]))
	case "drop":
		runDrop(m, len(rest) == 1)
	case "status":
		runStatus(m, src)
	}
}

func atoi(cmd, value string) int {
	n, err := strconv.Atoi(value)
	if err != nil {
		usageError(cmd, fmt.Sprintf("invalid value %q", value))
	}
	return n
}

func openMigrate() (*migrate.Migrate, source.Driver, error) {
	cfg := config.Setup()

	src, err := iofs.New(postgres.MigrationsFS, "migrations")
	if err != nil {
		return nil, nil, fmt.Errorf("load migration files: %w", err)
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
		return nil, nil, fmt.Errorf("initialize migrations: %w", err)
	}

	return m, src, nil
}

func usage(w io.Writer) {
	fmt.Fprint(w, usageText)
}

func usageError(cmd, msg string) {
	log.Errorf("%s: %s", cmd, msg)
	usage(os.Stderr)
	os.Exit(2)
}
