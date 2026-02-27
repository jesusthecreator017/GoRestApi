package store

import (
	"context"
	"fmt"

	"github.com/jesusthecreator017/fswithgo/internal/store/dbsqlc"
)

type IssueStatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type AdminStats struct {
	TotalUsers    int64              `json:"total_users"`
	TotalIssues   int64              `json:"total_issues"`
	IssuesByStatus []IssueStatusCount `json:"issues_by_status"`
}

type AdminStore struct {
	queries *dbsqlc.Queries
}

func (a *AdminStore) GetStats(ctx context.Context) (*AdminStats, error) {
	totalUsers, err := a.queries.CountUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("counting users: %w", err)
	}

	rows, err := a.queries.CountIssuesByStatus(ctx)
	if err != nil {
		return nil, fmt.Errorf("counting issues by status: %w", err)
	}

	var totalIssues int64
	issuesByStatus := make([]IssueStatusCount, len(rows))
	for i, row := range rows {
		issuesByStatus[i] = IssueStatusCount{
			Status: row.Status,
			Count:  row.Count,
		}
		totalIssues += row.Count
	}

	return &AdminStats{
		TotalUsers:     totalUsers,
		TotalIssues:    totalIssues,
		IssuesByStatus: issuesByStatus,
	}, nil
}
