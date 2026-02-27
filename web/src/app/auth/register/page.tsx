"use client";

import { useRegister } from "@/hooks/useAuth";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Card,
	CardAction,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";

export default function RegisterPage() {
	const [name, setName] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const { user } = useAuth();
	const router = useRouter();
	const register = useRegister();

	// Redirect if already authenticated
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
		register.mutate({ name, email, password });
	}

	return (
		<div className="flex min-h-screen items-center justify-center">
			<Card className="w-full max-w-sm">
				<CardHeader>
					<CardTitle>Register your account</CardTitle>
					<CardDescription>
						Enter your details below to create an account
					</CardDescription>
				</CardHeader>
				<CardContent>
					<form
						id="register-form"
						onSubmit={handleSubmit}
						className="flex flex-col gap-4"
					>
						<div className="flex flex-col gap-2">
							<Label htmlFor="name">Name</Label>
							<Input
								id="name"
								type="text"
								placeholder="Your name"
								value={name}
								onChange={(e) => setName(e.target.value)}
								required
							/>
						</div>
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
						{register.error && (
							<p className="text-destructive text-sm">
								{register.error.message}
							</p>
						)}
						<Button
							type="submit"
							className="w-full"
							disabled={register.isPending}
						>
							{register.isPending ? "Creating account..." : "Register"}
						</Button>
					</form>
				</CardContent>
				<CardFooter className="justify-center">
					<p className="text-muted-foreground text-sm">
						Already have an account?{" "}
						<Link
							href="/auth/login"
							className="text-primary underline underline-offset-4"
						>
							Login
						</Link>
					</p>
				</CardFooter>
			</Card>
		</div>
	);
}
