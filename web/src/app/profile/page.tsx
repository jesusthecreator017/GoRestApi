"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Mail, Calendar, Shield, User as UserIcon } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { useIssues } from "@/hooks/useIssues";
import {
	hasPermission,
	PERM_READ,
	PERM_WRITE,
	PERM_ADMIN,
} from "@/schemas/user";
import { UserAvatar } from "@/components/UserAvatar";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";

function formatDate(dateStr: string) {
	return new Date(dateStr).toLocaleDateString("en-US", {
		year: "numeric",
		month: "long",
		day: "numeric",
	});
}

export default function ProfilePage() {
	const { user, isLoading: authLoading } = useAuth();
	const router = useRouter();
	const { data: issuesData } = useIssues();

	useEffect(() => {
		if (!authLoading && !user) {
			router.push("/auth/login");
		}
	}, [authLoading, user, router]);

	if (authLoading || !user) {
		return null;
	}

	const userIssues =
		issuesData?.filter((issue) => issue.user_id === user.id) ?? [];
	const issuesByStatus = {
		Incomplete: userIssues.filter((i) => i.status === "Incomplete").length,
		"In-Progress": userIssues.filter((i) => i.status === "In-Progress").length,
		Complete: userIssues.filter((i) => i.status === "Complete").length,
	};

	const permissions: string[] = [];
	if (hasPermission(user.permissions, PERM_READ)) permissions.push("Read");
	if (hasPermission(user.permissions, PERM_WRITE)) permissions.push("Write");
	if (hasPermission(user.permissions, PERM_ADMIN)) permissions.push("Admin");

	return (
		<div className="mx-auto max-w-2xl px-4 py-8">
			<h1 className="mb-6 text-2xl font-bold">Profile</h1>

			<Card>
				<CardHeader>
					<div className="flex items-center gap-4">
						<UserAvatar name={user.name} className="size-16 text-xl" />
						<div>
							<CardTitle className="text-xl">{user.name}</CardTitle>
							<p className="text-muted-foreground text-sm">{user.email}</p>
						</div>
					</div>
				</CardHeader>
				<CardContent className="space-y-4">
					<Separator />

					<div className="grid gap-4 sm:grid-cols-2">
						<div className="flex items-center gap-2">
							<UserIcon className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Name</span>
						</div>
						<span className="text-sm font-medium">{user.name}</span>

						<div className="flex items-center gap-2">
							<Mail className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Email</span>
						</div>
						<span className="text-sm font-medium">{user.email}</span>

						<div className="flex items-center gap-2">
							<Shield className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">
								Permissions
							</span>
						</div>
						<div className="flex gap-1">
							{permissions.map((perm) => (
								<Badge key={perm} variant="secondary">
									{perm}
								</Badge>
							))}
						</div>

						<div className="flex items-center gap-2">
							<Calendar className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">Joined</span>
						</div>
						<span className="text-sm font-medium">
							{formatDate(user.created_at)}
						</span>

						<div className="flex items-center gap-2">
							<Calendar className="text-muted-foreground size-4" />
							<span className="text-muted-foreground text-sm">
								Last updated
							</span>
						</div>
						<span className="text-sm font-medium">
							{formatDate(user.updated_at)}
						</span>
					</div>
				</CardContent>
			</Card>

			<h2 className="mt-8 mb-4 text-lg font-semibold">Your Issues</h2>
			<div className="grid grid-cols-3 gap-4">
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByStatus.Incomplete}</p>
						<p className="text-muted-foreground text-sm">Incomplete</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">
							{issuesByStatus["In-Progress"]}
						</p>
						<p className="text-muted-foreground text-sm">In Progress</p>
					</CardContent>
				</Card>
				<Card>
					<CardContent className="pt-6 text-center">
						<p className="text-3xl font-bold">{issuesByStatus.Complete}</p>
						<p className="text-muted-foreground text-sm">Complete</p>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
