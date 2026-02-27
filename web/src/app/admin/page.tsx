"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { useAdminStats } from "@/hooks/useAdmin";
import { hasPermission, PERM_ADMIN } from "@/schemas/user";
import {
    Card,
    CardHeader,
    CardTitle,
    CardContent,
} from "@/components/ui/card";

export default function AdminPage() {
    const { user, isLoading: authLoading } = useAuth();
    const router = useRouter();

    const isAdmin = user ? hasPermission(user.permissions, PERM_ADMIN) : false;

    useEffect(() => {
        if (!authLoading && !isAdmin) {
            router.push("/");
        }
    }, [authLoading, isAdmin, router]);

    const { data, isLoading, error } = useAdminStats();

    if (authLoading || !isAdmin) {
        return null;
    }

    if (isLoading) {
        return (
            <div className="mx-auto max-w-4xl px-4 py-8">
                <p className="text-muted-foreground">Loading stats...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="mx-auto max-w-4xl px-4 py-8">
                <p className="text-destructive">Failed to load admin stats.</p>
            </div>
        );
    }

    const stats = data?.stats;

    return (
        <div className="mx-auto max-w-4xl px-4 py-8">
            <h1 className="mb-6 text-2xl font-bold">Admin Dashboard</h1>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                <Card>
                    <CardHeader>
                        <CardTitle>Total Users</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-3xl font-bold">{stats?.total_users ?? 0}</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Total Issues</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-3xl font-bold">{stats?.total_issues ?? 0}</p>
                    </CardContent>
                </Card>
                {stats?.issues_by_status?.map((item) => (
                    <Card key={item.status}>
                        <CardHeader>
                            <CardTitle>{item.status}</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <p className="text-3xl font-bold">{item.count}</p>
                        </CardContent>
                    </Card>
                ))}
            </div>
        </div>
    );
}
