import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function LandingPage() {
	return (
		<div className="flex min-h-[calc(100vh-3.5rem)] flex-col items-center justify-center gap-6">
			<h1 className="text-4xl font-bold">Track your issues</h1>
			<p className="text-muted-foreground text-lg">
				A simple issue tracker built with Go and Next.js
			</p>
			<div className="flex gap-3">
				<Button asChild>
					<Link href="/auth/login">Login</Link>
				</Button>
				<Button asChild variant="secondary">
					<Link href="/auth/register">Register</Link>
				</Button>
			</div>
		</div>
	);
}
