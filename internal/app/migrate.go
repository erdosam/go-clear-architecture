//go:build migrate

package app

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"log"
	"os"
	"strconv"
	"time"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func Migrate() {
}

func init() {
	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: PG_URL")
	}
	var (
		attempts  = _defaultAttempts
		err       error
		m         *migrate.Migrate
		srcDriver source.Driver
	)
	for attempts > 0 {
		migrationSrc, ok := os.LookupEnv("MIGRATION_SRC")
		if !ok {
			log.Fatalf("migrate: environment variable not declared: MIGRATION_SRC")
		}
		srcDriver, err = source.Open(migrationSrc)
		if err != nil {
			log.Fatalf(err.Error())
		}
		m, err = migrate.NewWithSourceInstance(migrationSrc, srcDriver, databaseURL)
		if err == nil {
			break
		}

		log.Printf("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	migrateMode := "up"
	s := 0
	if len(os.Args) > 1 {
		s, _ = strconv.Atoi(os.Args[1])
	}
	if s == 0 {
		err = m.Up()
	} else {
		err = m.Steps(s)
		if s < 0 {
			migrateMode = "down"
		}
	}
	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			log.Printf("Migration closed")
		}
	}(m)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: %s error: %s", migrateMode, err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}
	log.Printf("Migrate: %s success", migrateMode)
}
