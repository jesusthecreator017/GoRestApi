"use client";

import { useLogin } from "@/hooks/useAuth";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

export default function LoginPage() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const { user } = useAuth();
	const router = useRouter();
	const login = useLogin();

	// Redirect is already authenticated
	useEffect(() => {
		if (user) {
			router.push("/");
		}
	}, [user, router]);

	if (user) {
		return null;
	}

	function handleSubmit(e: React.SyntheticEvent<HTMLFormElement, SubmitEvent>) {
		e.preventDefault();
		login.mutate({ email, password });
	}

	return (
		<div className="flex min-h-screen items-center justify-center">
			<Card className="w-full max-w-sm">
				<CardHeader>
					<CardTitle>Login your account</CardTitle>
					<CardDescription>
						Enter your details below to login to your account
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form
						id="login-form"
						onSubmit={handleSubmit}
						className="flex flex-col gap-4"
					>
						<div className="flex flex-col gap-2">
							<Label htmlFor="email">Email</Label>
							<Input
								id="email"
								type="email"
								placeholder="you@example.com"
								value={email}
								onChange={(e) => setEmail(e.target.value)}
								required
							/>
						</div>
						<div className="flex flex-col gap-2">
							<Label htmlFor="password">Password</Label>
							<Input
								id="password"
								type="password"
								placeholder="Min 8 characters"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								required
								minLength={8}
							/>
						</div>
						{login.error && (
							<p className="text-destructive text-sm">{login.error.message}</p>
						)}
						<Button type="submit" className="w-full" disabled={login.isPending}>
							{login.isPending ? "Creating account..." : "Login"}
						</Button>
					</form>
				</CardContent>
				<CardFooter className="justify-center">
					<p className="text-muted-foreground text-sm">
						Don't have an account?{" "}
						<Link
							href="/auth/register"
							className="text-primary underline underline-offset-4"
						>
							Register
						</Link>
					</p>
				</CardFooter>
			</Card>
		</div>
	);
}
