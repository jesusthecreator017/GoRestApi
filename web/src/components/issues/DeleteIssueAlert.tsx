"use client";

import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useDeleteIssue } from "@/hooks/useIssues";

interface DeleteIssueAlertProps {
	issueId: number;
	issueTitle: string;
	open: boolean;
	onOpenChange: (open: boolean) => void;
}

export function DeleteIssueAlert({
	issueId,
	issueTitle,
	open,
	onOpenChange,
}: DeleteIssueAlertProps) {
	const deleteIssue = useDeleteIssue();

	function handleConfirm() {
		deleteIssue.mutate(issueId, {
			onSettled: () => onOpenChange(false),
		});
	}

	return (
		<AlertDialog open={open} onOpenChange={onOpenChange}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Delete issue</AlertDialogTitle>
					<AlertDialogDescription>
						Are you sure you want to delete "{issueTitle}"? This action cannot
						be undone.
					</AlertDialogDescription>
				</AlertDialogHeader>
				<AlertDialogFooter>
					<AlertDialogCancel>Cancel</AlertDialogCancel>
					<AlertDialogAction
						variant="destructive"
						onClick={handleConfirm}
						disabled={deleteIssue.isPending}
					>
						{deleteIssue.isPending ? "Deleting..." : "Delete"}
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	);
}
