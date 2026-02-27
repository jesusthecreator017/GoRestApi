import { apiFetch } from "./client";

export interface IssueStatusCount {
    status: string;
    count: number;
}

export interface AdminStats {
    total_users: number;
    total_issues: number;
    issues_by_status: IssueStatusCount[];
}

export const adminApi = {
    getStats: () =>
        apiFetch<{ stats: AdminStats }>("/v1/admin/stats"),
};
