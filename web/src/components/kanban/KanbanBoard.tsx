"use client";

import { useState } from "react";
import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useIssues } from "@/hooks/useIssues";
import { KanbanColumn } from "./KanbanColumn";
import { CreateIssueDialog } from "@/components/issues/CreateIssueDialog";
import type { StatusType } from "@/schemas/issue";

const columns: { status: StatusType; label: string }[] = [
	{ status: "Incomplete", label: "Incomplete" },
	{ status: "In-Progress", label: "In Progress" },
	{ status: "Complete", label: "Complete" },
];

export function KanbanBoard() {
	const [createOpen, setCreateOpen] = useState(false);
	const { data: issues, isLoading, error } = useIssues();

	if (isLoading) {
		return (
			<div className="flex items-center justify-center py-20">
				<p className="text-muted-foreground">Loading issues...</p>
			</div>
		);
	}

	if (error) {
		return (
			<div className="flex items-center justify-center py-20">
				<p className="text-destructive">
					Failed to load issues: {error.message}
				</p>
			</div>
		);
	}

	return (
		<div className="grid gap-6">
			<div className="flex items-center justify-between">
				<h1 className="text-2xl font-bold">Issues</h1>
				<Button onClick={() => setCreateOpen(true)}>
					<Plus className="size-4 mr-1" />
					New Issue
				</Button>
			</div>
			<div className="flex gap-4 overflow-x-auto pb-4">
				{columns.map((col) => (
					<KanbanColumn
						key={col.status}
						label={col.label}
						issues={(issues ?? []).filter((i) => i.status === col.status)}
					/>
				))}
			</div>
			<CreateIssueDialog open={createOpen} onOpenChange={setCreateOpen} />
		</div>
	);
}
