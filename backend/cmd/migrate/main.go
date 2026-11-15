// Command migrate applies the SQL migrations in ./migrations.
//
// The service is deployed as a container, so migrations run as a separate
// step (Render's pre-deploy command) instead of at application startup.
package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	dir := os.Getenv("MIGRATIONS_DIR")
	if dir == "" {
		dir = "migrations"
	}

	m, err := migrate.New("file://"+dir, withPgxScheme(databaseURL))
	if err != nil {
		log.Fatal("failed to open migrations: ", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("failed to apply migrations: ", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		log.Fatal("failed to read migration version: ", err)
	}

	log.Printf("database is at version %d (dirty=%t)", version, dirty)
}

// withPgxScheme points the URL at the pgx driver instead of lib/pq.
func withPgxScheme(databaseURL string) string {
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if strings.HasPrefix(databaseURL, prefix) {
			return "pgx5://" + strings.TrimPrefix(databaseURL, prefix)
		}
	}
	return databaseURL
}
