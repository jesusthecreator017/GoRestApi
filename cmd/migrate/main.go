package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	// Blank imports register drivers via their init() functions.
	// This is a common Go pattern — you never call these packages directly,
	// but importing them makes them available when migrate.New() looks up
	// the "pgx5://" scheme and "file://" source.
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// This program runs database migrations.
//
// Usage:
//
//	DATABASE_URL="pgx5://user:pass@localhost:5432/fswithgo?sslmode=disable" go run ./cmd/migrate -direction up
//	DATABASE_URL="pgx5://user:pass@localhost:5432/fswithgo?sslmode=disable" go run ./cmd/migrate -direction down
//
// IMPORTANT: The DATABASE_URL scheme must be "pgx5://" (not "postgres://")
// because golang-migrate's pgx/v5 driver registers itself under that scheme.
func main() {
	direction := flag.String("direction", "up", "migration direction (up or down)")
	flag.Parse()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// "file://..." points to the directory containing your .sql migration files.
	// migrate.New connects to the database and reads the migrations directory.
	m, err := migrate.New("file://cmd/migrate/migrations", dbURL)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		// m.Up() applies all pending migrations in order (000001, 000002, ...).
		// migrate.ErrNoChange means all migrations are already applied — not an error.
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failed: %v", err)
		}
		fmt.Println("migrations applied successfully")
	case "down":
		// m.Down() rolls back ALL migrations. In production you'd typically
		// use m.Steps(-1) to roll back just one migration at a time.
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down failed: %v", err)
		}
		fmt.Println("migrations rolled back successfully")
	default:
		log.Fatalf("unknown direction: %s (use 'up' or 'down')", *direction)
	}
}
