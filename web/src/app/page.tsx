"use client";

import { useAuth } from "@/contexts/AuthContext";
import Dashboard from "./Dashboard";
import LandingPage from "./LandingPage";

export default function HomePage() {
	const { user, isLoading } = useAuth();

	if (isLoading) {
		return null;
	}

	if (!user) {
		return <LandingPage />;
	}

	return (
		<div className="mx-auto max-w-6xl px-4 py-8">
			<Dashboard />
		</div>
	);
}
