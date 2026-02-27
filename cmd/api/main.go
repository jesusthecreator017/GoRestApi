package main

import (
	"log"
	"os"
	"time"

	"github.com/jesusthecreator017/fswithgo/internal/db"
	"github.com/jesusthecreator017/fswithgo/internal/env"
	"github.com/jesusthecreator017/fswithgo/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":3000"),
		db: dbConfig{
			// DATABASE_URL is the standard env var for PostgreSQL connection strings.
			// Format: "postgres://user:password@host:port/dbname?sslmode=disable"
			dsn:         env.GetString("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/fswithgo?sslmode=disable"),
			maxConns:    env.GetInt("DB_MAX_CONNS", 25),
			maxIdleTime: time.Duration(env.GetInt("DB_MAX_IDLE_MINS", 15)) * time.Minute,
		},
		corsOrigin: os.Getenv("CORS_ORIGIN"),
		jwtSecret:  env.GetString("JWT_SECRET", ""),
	}

	// Initialize any environment variables

	// Create the pgxpool connection pool.
	// This is like *sql.DB but for pgx — it manages a pool of connections
	// so you don't open/close a connection for every query.
	// db.New() also pings the database to verify it's reachable.
	pool, err := db.New(cfg.db.dsn, cfg.db.maxConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// pool.Close() returns all connections to PostgreSQL when the server shuts down.
	defer pool.Close()
	log.Println("database connection pool established")

	// Pass the pool to NewStorage. Inside, it creates a sqlc Queries struct
	// (dbsqlc.New(pool)) and wires it into each repository implementation.
	store := store.NewStorage(pool)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
