package store

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

// Storage holds all repository interfaces — one per domain entity.
//
// Each field is an interface (not a concrete type) so you can:
//   - Swap in mock implementations for unit testing
//   - Change the underlying database without touching handler code
//
// This is the "Repository Pattern" — your HTTP handlers depend on these
// interfaces, not on the database directly.
type Storage struct {
	Issues interface {
		Create(context.Context, *Issue) error
		GetByID(context.Context, int64) (*Issue, error)
		List(context.Context) ([]*Issue, error)
		ListByUserID(context.Context, uuid.UUID) ([]*Issue, error)
		UpdateStatus(context.Context, int64, StatusType) (*Issue, error)
		Delete(context.Context, int64) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByEmail(context.Context, string) (*User, error)
		GetByID(context.Context, uuid.UUID) (*User, error)
	}
	Admin interface {
		GetStats(context.Context) (*AdminStats, error)
	}
}

// NewStorage creates a Storage with real database-backed implementations.
//
// It takes a *pgxpool.Pool (not *sql.DB) because we use pgx/v5 native mode.
// dbsqlc.New(pool) creates the sqlc-generated Queries struct, which accepts
// anything implementing the DBTX interface — pgxpool.Pool satisfies it,
// and so do pgx.Conn and pgx.Tx (for transactions).
func NewStorage(pool *pgxpool.Pool) Storage {
	queries := dbsqlc.New(pool)

	return Storage{
		Issues: &IssueStore{queries: queries},
		Users:  &UserStore{queries: queries},
		Admin:  &AdminStore{queries: queries},
	}
}
