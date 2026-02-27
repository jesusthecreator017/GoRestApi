"use client";

import { useState } from "react";
import Link from "next/link";
import { Trash2 } from "lucide-react";
import {
	Card,
	CardContent,
	CardHeader,
	CardTitle,
	CardAction,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { useUpdateIssueStatus } from "@/hooks/useIssues";
import { DeleteIssueAlert } from "@/components/issues/DeleteIssueAlert";
import type { Issue, StatusType } from "@/schemas/issue";

const statusColors: Record<StatusType, string> = {
	Incomplete: "bg-red-100 text-red-800 border-red-200",
	"In-Progress": "bg-yellow-100 text-yellow-800 border-yellow-200",
	Complete: "bg-green-100 text-green-800 border-green-200",
};

export function IssueCard({ issue }: { issue: Issue }) {
	const [deleteOpen, setDeleteOpen] = useState(false);
	const updateStatus = useUpdateIssueStatus();

	function handleStatusChange(status: string) {
		updateStatus.mutate({ id: issue.id, status: status as StatusType });
	}

	return (
		<>
			<Card className="gap-3 py-4">
				<CardHeader className="gap-1">
					<CardTitle className="text-sm">
						<Link
							href={`/issues/${issue.id}`}
							className="text-foreground no-underline hover:underline"
						>
							{issue.title}
						</Link>
					</CardTitle>
					<CardAction>
						<Button
							variant="ghost"
							size="icon-xs"
							onClick={() => setDeleteOpen(true)}
						>
							<Trash2 className="size-3.5 text-muted-foreground" />
						</Button>
					</CardAction>
				</CardHeader>
				<CardContent className="grid gap-3">
					{issue.description && (
						<p className="text-muted-foreground text-xs line-clamp-2">
							{issue.description}
						</p>
					)}
					<div className="flex items-center justify-between gap-2">
						<Badge className={statusColors[issue.status]} variant="outline">
							{issue.status}
						</Badge>
						<span className="text-muted-foreground text-xs truncate">
							{issue.user_name}
						</span>
					</div>
					<Select value={issue.status} onValueChange={handleStatusChange}>
						<SelectTrigger size="sm" className="w-full text-xs text-white">
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							<SelectItem value="Incomplete">Incomplete</SelectItem>
							<SelectItem value="In-Progress">In-Progress</SelectItem>
							<SelectItem value="Complete">Complete</SelectItem>
						</SelectContent>
					</Select>
				</CardContent>
			</Card>
			<DeleteIssueAlert
				issueId={issue.id}
				issueTitle={issue.title}
				open={deleteOpen}
				onOpenChange={setDeleteOpen}
			/>
		</>
	);
}
