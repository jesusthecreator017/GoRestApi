package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// New creates a PostgreSQL connection pool using pgx/v5 native mode.
//
// pgxpool.Pool is the pgx equivalent of *sql.DB — it manages a pool of
// database connections. Unlike database/sql, pgx talks directly to PostgreSQL
// using its native binary protocol, which is faster and supports PostgreSQL-
// specific types (UUID, arrays, JSONB, etc.) without conversion overhead.
//
// The pool opens connections lazily — it doesn't connect until a query runs.
// We call Ping() at the end to verify the connection works immediately.
func New(dsn string, maxConns int, maxIdleTime time.Duration) (*pgxpool.Pool, error) {
	// ParseConfig reads a PostgreSQL connection string (DSN) and returns
	// a pool configuration. Example DSN:
	//   "postgres://user:pass@localhost:5432/fswithgo?sslmode=disable"
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing database config: %w", err)
	}

	// MaxConns limits how many connections the pool will open at once.
	// If all connections are in use, new queries wait until one is returned.
	config.MaxConns = int32(maxConns)

	// MaxConnIdleTime closes connections that have been sitting idle
	// for this long, freeing resources on both Go and PostgreSQL sides.
	config.MaxConnIdleTime = maxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	// Ping sends a simple query to verify the database is reachable.
	// Without this, you wouldn't know about connection problems until
	// the first real query, which makes debugging harder.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return pool, nil
}
