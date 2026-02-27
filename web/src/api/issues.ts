import { apiFetch } from "./client";
import type { Issue, CreateIssueInput, StatusType } from "../schemas/issue";

export const issuesApi = {
    list: () =>
        apiFetch<{ issues: Issue[] }>("/v1/issues").then((r) => r.issues),

    getById: (id: number) =>
        apiFetch<{ issue: Issue }>(`/v1/issues/${id}`).then((r) => r.issue),

    create: (data: CreateIssueInput) =>
        apiFetch<{ issue: Issue }>("/v1/issues", {
            method: "POST",
            body: JSON.stringify(data),
        }).then((r) => r.issue),

    updateStatus: (id: number, status: StatusType) =>
        apiFetch<{ issue: Issue }>(`/v1/issues/${id}/status`, {
            method: "PATCH",
            body: JSON.stringify({ status }),
        }).then((r) => r.issue),

    delete: (id: number) =>
        apiFetch<{ message: string }>(`/v1/issues/${id}`, {
            method: "DELETE"
        }),
}