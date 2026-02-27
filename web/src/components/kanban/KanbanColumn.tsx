"use client";

import { Badge } from "@/components/ui/badge";
import { IssueCard } from "./IssueCard";
import type { Issue } from "@/schemas/issue";

interface KanbanColumnProps {
  label: string;
  issues: Issue[];
}

export function KanbanColumn({ label, issues }: KanbanColumnProps) {
  return (
    <div className="flex flex-1 flex-col gap-3 min-w-70">
      <div className="flex items-center gap-2 px-1">
        <h2 className="text-sm font-semibold">{label}</h2>
        <Badge variant="secondary">{issues.length}</Badge>
      </div>
      <div className="flex flex-col gap-2">
        {issues.map((issue) => (
          <IssueCard key={issue.id} issue={issue} />
        ))}
        {issues.length === 0 && (
          <p className="text-muted-foreground text-center text-sm py-8">
            No issues
          </p>
        )}
      </div>
    </div>
  );
}
