"use client";

import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";
import { issuesApi } from "@/api/issues";

export default function IssuePage() {
  const { issueId } = useParams<{ issueId: string }>();
  const id = Number(issueId);

  const { data: issue, isLoading, error } = useQuery({
    queryKey: ["issues", id],
    queryFn: () => issuesApi.getById(id),
    enabled: !isNaN(id),
  });

  if (isLoading) {
    return <div className="mx-auto max-w-2xl px-4 py-8">Loading...</div>;
  }

  if (error || !issue) {
    return (
      <div className="mx-auto max-w-2xl px-4 py-8">
        <p className="text-destructive">
          Failed to load issue{error ? `: ${error.message}` : ""}
        </p>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-2xl font-bold">{issue.title}</h1>
      <p className="mt-4 text-muted-foreground">{issue.description}</p>
    </div>
  );
}
