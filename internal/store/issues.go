package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

type StatusType string

const (
	Incomplete StatusType = "Incomplete"
	InProgress StatusType = "In-Progress"
	Complete   StatusType = "Complete"
)

type Issue struct {
	ID          int64      `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	UserName    string     `json:"user_name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      StatusType `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type IssueStore struct {
	queries *dbsqlc.Queries
}

var ErrNotFound = errors.New("resource not found")

func (s *IssueStore) Create(ctx context.Context, issue *Issue) error {
	row, err := s.queries.CreateIssue(ctx, dbsqlc.CreateIssueParams{
		Title:       issue.Title,
		UserID:      issue.UserID,
		Description: issue.Description,
		Status:      string(issue.Status),
	})
	if err != nil {
		return fmt.Errorf("creating issue: %w", err)
	}

	issue.ID = row.ID
	issue.CreatedAt = row.CreatedAt.Time
	issue.UpdatedAt = row.UpdatedAt.Time
	return nil
}

func (s *IssueStore) GetByID(ctx context.Context, id int64) (*Issue, error) {
	row, err := s.queries.GetIssueByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("getting issue: %w", err)
	}

	return issueRowToDomain(row), nil
}

func (s *IssueStore) List(ctx context.Context) ([]*Issue, error) {
	rows, err := s.queries.ListIssues(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing issues: %w", err)
	}

	issues := make([]*Issue, len(rows))
	for i, row := range rows {
		issues[i] = listRowToDomain(row)
	}
	return issues, nil
}

func (s *IssueStore) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Issue, error) {
	rows, err := s.queries.ListIssuesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing issues by user: %w", err)
	}

	issues := make([]*Issue, len(rows))
	for i, row := range rows {
		issues[i] = &Issue{
			ID:          row.ID,
			UserID:      row.UserID,
			UserName:    row.UserName,
			Title:       row.Title,
			Description: row.Description,
			Status:      StatusType(row.Status),
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}
	}
	return issues, nil
}

func (s *IssueStore) UpdateStatus(ctx context.Context, id int64, newStatus StatusType) (*Issue, error) {
	row, err := s.queries.UpdateIssueStatus(ctx, dbsqlc.UpdateIssueStatusParams{
		ID:     id,
		Status: string(newStatus),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("updating issue status: %w", err)
	}

	return toDomainIssue(row), nil
}

func (s *IssueStore) Delete(ctx context.Context, id int64) error {
	err := s.queries.DeleteIssue(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("deleting issue: %w", err)
	}
	return nil
}

// toDomainIssue converts a sqlc-generated Issue (from CreateIssue/UpdateIssueStatus)
// to the domain Issue. These queries don't JOIN users, so UserName is left empty.
func toDomainIssue(row dbsqlc.Issue) *Issue {
	return &Issue{
		ID:          row.ID,
		UserID:      row.UserID,
		Title:       row.Title,
		Description: row.Description,
		Status:      StatusType(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

// issueRowToDomain converts a GetIssueByIDRow (which includes user_name via JOIN).
func issueRowToDomain(row dbsqlc.GetIssueByIDRow) *Issue {
	return &Issue{
		ID:          row.ID,
		UserID:      row.UserID,
		UserName:    row.UserName,
		Title:       row.Title,
		Description: row.Description,
		Status:      StatusType(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}

// listRowToDomain converts a ListIssuesRow (which includes user_name via JOIN).
func listRowToDomain(row dbsqlc.ListIssuesRow) *Issue {
	return &Issue{
		ID:          row.ID,
		UserID:      row.UserID,
		UserName:    row.UserName,
		Title:       row.Title,
		Description: row.Description,
		Status:      StatusType(row.Status),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
