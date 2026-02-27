"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { issuesApi } from "@/api/issues";
import type { CreateIssueInput, StatusType } from "@/schemas/issue";

export function useIssues() {
    return useQuery({
        queryKey: ["issues"],
        queryFn: issuesApi.list
    });
}

export function useCreateIssue() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (data: CreateIssueInput) => issuesApi.create(data),
        onSuccess: () => qc.invalidateQueries({ queryKey: ["issues"] }),
    });
}

export function useUpdateIssueStatus() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: ({ id, status }: { id: number; status: StatusType }) =>
            issuesApi.updateStatus(id, status),
        onSuccess: () => qc.invalidateQueries({ queryKey: ["issues"] }),
    });
}

export function useDeleteIssue() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (id: number) => issuesApi.delete(id),
        onSuccess: () => qc.invalidateQueries({ queryKey: ["issues"] }),
    });
}
