"use client";

import { useAuth } from "@/contexts/AuthContext";

export default function Dashboard() {
	const { user } = useAuth();

	return (
		<div>
			<h1 className="text-2xl font-bold">Welcome, {user?.name}</h1>
		</div>
	);
}
