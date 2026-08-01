package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
)

func runUp(m *migrate.Migrate) {
	err := m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("No pending migrations")
		return
	}
	if err != nil {
		log.Fatal("Failed to run migrations!", "err", err.Error())
	}
	log.Info("Database migrations completed successfully")
}

func runDown(m *migrate.Migrate) {
	err := m.Steps(-1)
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("No migrations to roll back")
		return
	}
	if err != nil {
		log.Fatal("Failed to roll back migration!", "err", err.Error())
	}
	log.Info("Database migration rolled back successfully")
}

func runSteps(m *migrate.Migrate, steps int) {
	err := m.Steps(steps)
	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("No change")
		return
	}
	if err != nil {
		log.Fatal("Failed to apply steps!", "err", err.Error())
	}
	log.Info("Migrations applied successfully")
}

func runVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		log.Info("No migrations applied yet")
		return
	}
	if err != nil {
		log.Fatal("Failed to get current version!", "err", err.Error())
	}
	log.Infof("Current version: %d (dirty: %t)", version, dirty)
}

func runForce(m *migrate.Migrate, version int) {
	if err := m.Force(version); err != nil {
		log.Fatal("Failed to force version!", "err", err.Error())
	}
	log.Infof("Forced version to %d", version)
}

func runDrop(m *migrate.Migrate, force bool) {
	if !force {
		fmt.Fprint(os.Stdout, "Are you sure you want to drop all database objects? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		answer, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Failed to read confirmation!", "err", err.Error())
		}
		if !strings.EqualFold(strings.TrimSpace(answer), "y") {
			log.Info("Drop aborted")
			return
		}
	}

	if err := m.Drop(); err != nil {
		log.Fatal("Failed to drop database objects!", "err", err.Error())
	}
	log.Info("Database objects dropped successfully")
}

func runStatus(m *migrate.Migrate, src source.Driver) {
	applied := false
	var current uint
	current, dirty, err := m.Version()
	switch {
	case errors.Is(err, migrate.ErrNilVersion):
		log.Info("No migrations applied yet")
	case err != nil:
		log.Fatal("Failed to get current version!", "err", err.Error())
	default:
		applied = true
		log.Infof("Current version: %d (dirty: %t)", current, dirty)
	}

	first, err := src.First()
	if err != nil {
		log.Fatal("Failed to read migrations!", "err", err.Error())
	}

	for version := first; ; {
		r, name, err := src.ReadUp(version)
		if err != nil {
			log.Fatal("Failed to read migration!", "err", err.Error())
		}
		r.Close()

		state := "pending"
		if applied && version <= current {
			state = "applied"
		}
		log.Infof("  %d %s [%s]", version, name, state)

		next, err := src.Next(version)
		if err != nil {
			break
		}
		version = next
	}
}
